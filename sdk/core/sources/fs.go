package sources

import (
	"context"
	"io/fs"

	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"go.uber.org/multierr"
)

// FsEntry represents a single entry (file or directory) within a filesystem.
//
// It includes the relative path of the entry within the fs.FS hierarchy and its
// associated fs.DirEntry metadata. FsEntry values are returned by filesystem-based
// core.Source implementations like NewFs and NewFsErr.
type FsEntry struct {
	Path     string      // Relative path of the entry within the fs.FS
	DirEntry fs.DirEntry // Metadata describing the entry (file type, name, etc.)
}

// NewFs constructs a new core.Source[FsEntry] that lists files and directories
// from a given fs.FS, applying a boolean predicate to determine which entries
// to include in the output.
//
// This function performs a recursive traversal of the provided fs.FS using
// fs.WalkDir starting at the root path ".". For each visited entry, it calls
// the supplied predicate function. If predicate(path, entry) returns true,
// the entry is included in the result.
//
// Parameters:
//
//	f         — an implementation of fs.FS (e.g., os.DirFS("."), embed.FS, etc.).
//	predicate — a function that determines whether to include a given entry. It
//	            receives the entry's relative path and its fs.DirEntry.
//
// Returns:
//
//	core.Source[FsEntry] — a composable data source of filesystem entries that
//	                       satisfy the predicate.
//
// Behavior:
//   - Recursively traverses the fs.FS starting from the root (".").
//   - The predicate is evaluated for every file and directory encountered.
//   - Only entries for which predicate(path, entry) returns true are included.
//   - If fs.WalkDir encounters an error, traversal stops and the error is returned.
//   - All returned paths are relative to the fs.FS root.
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
//	    // Define a filesystem source that includes only Go source files.
//	    src := sources.NewFs(os.DirFS("."), func(path string, entry fs.DirEntry) bool {
//	        return !entry.IsDir() && strings.HasSuffix(path, ".go")
//	    })
//
//	    // Load matching entries from the source.
//	    entries, err := src.Load(context.Background())
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    for _, e := range entries {
//	        fmt.Println("Go file:", e.Path)
//	    }
//	}
func NewFs(f fs.FS, predicate func(string, fs.DirEntry) bool) core.Source[FsEntry] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]FsEntry, error) {
		var out []FsEntry

		// Recursively walk the filesystem starting at the fs.FS root (".").
		err := fs.WalkDir(f, ".", func(path string, entry fs.DirEntry, walkErr error) error {
			// Apply the predicate — include only matching entries.
			if predicate(path, entry) {
				out = append(out, FsEntry{
					Path:     path,
					DirEntry: entry,
				})
			}
			// Propagate traversal errors (e.g., permission denied, fs.SkipDir).
			return walkErr
		})

		// Return the collected entries and any traversal error.
		return out, err
	})
}

// NewFsErr constructs a new core.Source[FsEntry] that traverses an fs.FS and
// filters entries using a predicate that may also return an error.
//
// This function behaves like NewFs but provides more flexibility: the predicate
// can signal both inclusion and diagnostic errors. All errors returned by the
// predicate and the traversal itself are aggregated using go.uber.org/multierr,
// ensuring traversal continues even when individual predicate calls fail.
//
// Parameters:
//
//	f         — an fs.FS implementation (e.g., os.DirFS("."), embed.FS, etc.).
//	predicate — a function called for each visited entry that returns:
//	             (include, err):
//	               • include = true → include entry in output
//	               • err != nil    → append to aggregated error list
//
// Returns:
//
//	core.Source[FsEntry] — a composable data source of entries that match the predicate.
//
// Behavior:
//   - Traverses the filesystem recursively from the root using fs.WalkDir.
//   - For each entry, predicate(path, entry) is invoked.
//   - Includes entries where include == true.
//   - Collects all predicate and traversal errors via multierr.Append.
//   - Continues traversal even when predicate returns an error.
//   - Returns both results and the aggregated error at the end.
//
// Example:
//
//	src := sources.NewFsErr(os.DirFS("."), func(path string, entry fs.DirEntry) (bool, error) {
//	    if entry.IsDir() {
//	        return false, nil // Skip directories
//	    }
//	    if strings.HasSuffix(path, ".yaml") {
//	        return true, nil // Include YAML files
//	    }
//	    return false, nil
//	})
//
//	files, err := src.Load(context.Background())
//	if err != nil {
//	    log.Println("Errors encountered:", err)
//	}
//	for _, e := range files {
//	    fmt.Println("YAML file:", e.Path)
//	}
func NewFsErr(f fs.FS, predicate func(string, fs.DirEntry) (bool, error)) core.Source[FsEntry] {
	return core.MakeSourceFunc(func(ctx context.Context) ([]FsEntry, error) {
		var out []FsEntry
		var errs error

		// WalkDir recursively visits each file and directory in the filesystem.
		err := fs.WalkDir(f, ".", func(path string, entry fs.DirEntry, walkErr error) error {
			// Evaluate predicate for current entry; may produce inclusion flag and/or error.
			ok, err := predicate(path, entry)

			// Aggregate predicate errors but continue traversal.
			if err != nil {
				errs = multierr.Append(errs, err)
			}

			// Include entry in output only if predicate returned ok == true.
			if ok {
				out = append(out, FsEntry{
					Path:     path,
					DirEntry: entry,
				})
			}

			// Return traversal error (if any), preserving WalkDir semantics (e.g., fs.SkipDir).
			return walkErr
		})

		// Merge traversal-level error with all predicate-level errors.
		errs = multierr.Append(errs, err)
		return out, errs
	})
}
