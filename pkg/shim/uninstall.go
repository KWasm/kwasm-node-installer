package shim

import (
	"fmt"
	"os"

	"github.com/kwasm/kwasm-node-installer/pkg/state"
)

func (c *Config) Uninstall(shimName string) (string, error) {

	st, err := state.Get(c.config)
	if err != nil {
		return "", err
	}
	s := st.Shims[shimName]
	if s == nil {
		return "", fmt.Errorf("shim '%s' not installed", shimName)
	}
	filePath := s.Path
	filePathHost := c.config.PathWithHost(filePath)

	err = os.Remove(filePathHost)
	if err != nil {
		return "", err
	}

	st.RemoveShim(shimName)
	if err := st.Write(); err != nil {
		return "", err
	}
	return filePath, err
}
