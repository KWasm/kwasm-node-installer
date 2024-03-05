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

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/kwasm/kwasm-node-installer/pkg/containerd"
	"github.com/kwasm/kwasm-node-installer/pkg/shim"
)

// uninstallCmd represents the uninstall command.
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall containerd shims",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUninstall(cmd, args); err != nil {
			slog.Error("failed to uninstall", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(_ *cobra.Command, _ []string) error {
	slog.Info("uninstall called")
	shimName := conf.runtime.name
	runtimeName := path.Join(conf.kwasm.path, "bin", shimName)

	rootFs := afero.NewOsFs()
	hostFs := afero.NewBasePathFs(rootFs, conf.host.rootPath)

	containerdConfig := containerd.NewConfig(hostFs, conf.runtime.configPath)
	shimConfig := shim.NewConfig(rootFs, hostFs, conf.kwasm.assetPath, conf.kwasm.path)

	binPath, err := shimConfig.Uninstall(shimName)
	if err != nil {
		return fmt.Errorf("failed to delete shim '%s': %w", runtimeName, err)
	}

	configChanged, err := containerdConfig.RemoveRuntime(binPath)
	if err != nil {
		return fmt.Errorf("failed to write conteainerd config for shim '%s': %w", runtimeName, err)
	}

	if !configChanged {
		slog.Info("nothing changed, nothing more to do")
		return nil
	}

	slog.Info("restarting containerd")
	err = containerdConfig.RestartRuntime()
	if err != nil {
		return fmt.Errorf("failed to restart containerd: %w", err)
	}

	return nil
}
