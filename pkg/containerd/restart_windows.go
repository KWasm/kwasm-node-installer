//go:build windows
// +build windows

package containerd

import "errors"

func (c *Config) RestartRuntime() error {
	return errors.New("restarting containerd not implemented")
}
