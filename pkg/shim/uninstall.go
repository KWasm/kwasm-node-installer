package shim

import (
	"fmt"

	"github.com/kwasm/kwasm-node-installer/pkg/state"
)

func (c *Config) Uninstall(shimName string) (string, error) {

	st, err := state.Get(c.hostFs, c.kwasmPath)
	if err != nil {
		return "", err
	}
	s := st.Shims[shimName]
	if s == nil {
		return "", fmt.Errorf("shim '%s' not installed", shimName)
	}
	filePath := s.Path

	err = c.hostFs.Remove(filePath)
	if err != nil {
		return "", err
	}

	st.RemoveShim(shimName)
	if err := st.Write(); err != nil {
		return "", err
	}
	return filePath, err
}
