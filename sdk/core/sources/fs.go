package sources

import (
	"context"
	"io/fs"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// NewFs creates a new core.Source[string] that enumerates files within an fs.FS
// filesystem based on a user-provided filtering predicate.
//
// The resulting source yields a list of file paths (as strings) that satisfy
// the given predicate. It leverages fs.WalkDir to traverse the filesystem
// starting from the root (".").
//
// Parameters:
//
//	f         — an fs.FS filesystem implementation (e.g., os.DirFS, embed.FS, etc.)
//	predicate — a filter function that determines whether a file or directory
//	            path should be included in the output. It receives both the
//	            path (relative to root) and the corresponding fs.DirEntry.
//
// Returns:
//
//	core.Source[string] — a data source that lists all matching file paths
//	error               — any error encountered during directory traversal
//
// Behavior:
//   - All paths are relative to the provided fs.FS root.
//   - If fs.WalkDir encounters an error, traversal stops and the error is returned.
//   - The predicate is applied to every visited entry (files and directories).
//   - Only paths for which predicate(path, entry) returns true are included.
//
// Example:
//
//	import (
//	    "context"
//	    "io/fs"
//	    "os"
//	    "strings"
//	    "github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
//	)
//
//	func main() {
//	    // Create a filesystem source rooted at the current directory.
//	    src := sources.NewFs(os.DirFS("."), func(path string, entry fs.DirEntry) bool {
//	        // Include only Go source files.
//	        return !entry.IsDir() && strings.HasSuffix(path, ".go")
//	    })
//
//	    // Load all matching file paths.
//	    files, err := src.Load(context.Background())
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    for _, f := range files {
//	        fmt.Println("Go file:", f)
//	    }
//	}
func NewFs(f fs.FS, predicate func(string, fs.DirEntry) bool) core.Source[string] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]string, error) {
		var out []string

		// Walk the filesystem rooted at "." (the fs.FS root)
		err := fs.WalkDir(f, ".", func(path string, entry fs.DirEntry, walkErr error) error {
			// Apply the predicate — include the path if it matches
			if predicate(path, entry) {
				out = append(out, path)
			}
			// Return the traversal error (if any) so WalkDir can handle it
			return walkErr
		})

		return out, err
	})
}

// NewFsErr creates a core.Source[string] that walks an fs.FS filesystem
// and filters entries using a predicate that can return an error.
//
// It recursively visits all files and directories under the root ("."),
// invoking the predicate for each entry. If the predicate returns an error,
// it is aggregated using go.uber.org/multierr but traversal continues.
//
// Parameters:
//
//	f         — an fs.FS filesystem (e.g., os.DirFS("."), embed.FS, etc.)
//	predicate — a function called for each visited path and entry:
//	             (include, error):
//	               • include = true → include path in output
//	               • error != nil   → append to aggregated error list
//
// Returns:
//
//	core.Source[string] — a composable data source of file paths
//
// Behavior:
//   - Walks the fs.FS recursively using fs.WalkDir starting at the root.
//   - Calls predicate(path, entry) for every visited item.
//   - Includes paths where predicate returns true.
//   - Aggregates all predicate and WalkDir errors via multierr.Append.
//   - Continues traversal even if predicate returns errors.
//
// Example:
//
//	src := sources.NewFsErr(os.DirFS("."), func(path string, entry fs.DirEntry) (bool, error) {
//	    if entry.IsDir() {
//	        return false, nil
//	    }
//	    if strings.HasSuffix(path, ".yaml") {
//	        return true, nil
//	    }
//	    return false, nil
//	})
//
//	files, err := src.Load(context.Background())
//	if err != nil {
//	    log.Println("Errors:", err)
//	}
//	fmt.Println("YAML files:", files)
func NewFsErr(f fs.FS, predicate func(string, fs.DirEntry) (bool, error)) core.Source[string] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]string, error) {
		var out []string
		var errs error

		// WalkDir visits each file/directory recursively starting at the fs root
		err := fs.WalkDir(f, ".", func(path string, entry fs.DirEntry, _err error) error {
			// Evaluate predicate — include only if ok == true
			ok, err := predicate(path, entry)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
			if ok {
				out = append(out, path)
			}
			return _err // propagate WalkDir-level control (e.g., fs.SkipDir)
		})

		// Combine WalkDir-level error (if any) with predicate-level ones
		errs = multierr.Append(errs, err)
		return out, errs
	})
}
