/*
   Copyright The KWasm Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package containerd

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/kwasm/kwasm-node-installer/pkg/shim"
	"github.com/spf13/afero"
)

type Config struct {
	hostFs     afero.Fs
	configPath string
}

func NewConfig(hostFs afero.Fs, configPath string) *Config {
	return &Config{
		hostFs:     hostFs,
		configPath: configPath,
	}
}

func (c *Config) AddRuntime(shimPath string) error {
	runtimeName := shim.RuntimeName(path.Base(shimPath))

	cfg := generateConfig(shimPath, runtimeName)

	// Containerd config file needs to exist, otherwise return the error
	data, err := afero.ReadFile(c.hostFs, c.configPath)
	if err != nil {
		return err
	}

	// Warn if config.toml already contains runtimeName
	if strings.Contains(string(data), runtimeName) {
		slog.Info(fmt.Sprintf("config for runtime '%s' already exists, skipping", runtimeName))
		return nil
	}

	// Open file in append mode
	file, err := c.hostFs.OpenFile(c.configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append config
	_, err = file.WriteString(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RemoveRuntime(shimPath string) error {
	runtimeName := shim.RuntimeName(path.Base(shimPath))

	cfg := generateConfig(shimPath, runtimeName)

	// Containerd config file needs to exist, otherwise return the error
	data, err := afero.ReadFile(c.hostFs, c.configPath)
	if err != nil {
		return err
	}

	// Warn if config.toml does not contain the runtimeName
	if !strings.Contains(string(data), runtimeName) {
		slog.Warn(fmt.Sprintf("config for runtime '%s' does not exist, skipping", runtimeName))
		return nil
	}

	// Convert the file data to a string and replace the target string with an empty string.
	modifiedData := strings.Replace(string(data), cfg, "", -1)

	// Write the modified data back to the file.
	err = afero.WriteFile(c.hostFs, c.configPath, []byte(modifiedData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateConfig(shimPath string, runtimeName string) string {
	return fmt.Sprintf(`
# KWASM runtime config for %s
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.%s]
runtime_type = "%s"
`, runtimeName, runtimeName, shimPath)
}
