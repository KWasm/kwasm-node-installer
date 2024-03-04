package tests

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func FixtureFs(fixturePath string) afero.Fs {
	baseFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join("../..", fixturePath))
	fs := afero.NewCopyOnWriteFs(baseFs, afero.NewMemMapFs())
	return fs
}
