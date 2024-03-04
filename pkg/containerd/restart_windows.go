//go:build windows
// +build windows

package containerd

func (c *Config) RestartRuntime() error {
	panic("restarting containerd not implemented")
}
