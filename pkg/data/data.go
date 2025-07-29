package data

import (
	"embed"
	"io/fs"
	"sync"
)

const (
	crdsFolder = "crds"
)

//go:embed crds
var crdsFs embed.FS

func crds() (fs.FS, error) {
	return sub(crdsFs, crdsFolder)
}

func sub(f embed.FS, dir string) (fs.FS, error) {
	return fs.Sub(f, dir)
}

var (
	Crds = sync.OnceValues(crds)
)
