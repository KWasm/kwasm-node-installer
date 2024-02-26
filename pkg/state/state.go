package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type state struct {
	Shims        map[string]*Shim `json:"shims"`
	fs           afero.Fs
	lockFilePath string
}

func Get(fs afero.Fs, kwasmPath string) (*state, error) {
	out := state{
		Shims:        make(map[string]*Shim),
		lockFilePath: filepath.Join(kwasmPath, "kwasm-lock.json"),
		fs:           fs,
	}
	content, err := afero.ReadFile(fs, out.lockFilePath)
	if err == nil {
		err := json.Unmarshal(content, &out)
		return &out, err
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return &out, nil
}

func (l *state) ShimChanged(shimName string, sha256 []byte, path string) bool {
	shim, ok := l.Shims[shimName]
	if !ok {
		return true
	}

	return !bytes.Equal(shim.Sha256, sha256) || shim.Path != path
}

func (l *state) UpdateShim(shimName string, shim Shim) {
	l.Shims[shimName] = &shim
}

func (l *state) RemoveShim(shimName string) {
	delete(l.Shims, shimName)
}

func (l *state) Write() error {
	out, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		return err
	}

	slog.Info("writing lock file", "content", string(out))

	return afero.WriteFile(l.fs, l.lockFilePath, out, 0644)
}