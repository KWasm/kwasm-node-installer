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
	"io/fs"
	"log/slog"
	"os"
	"path"

	"github.com/kwasm/kwasm-node-installer/pkg/containerd"
	"github.com/kwasm/kwasm-node-installer/pkg/shim"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// installCmd represents the install command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install containerd shims",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInstall(cmd, args); err != nil {
			slog.Error("failed to install", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	installCmd.Flags().StringVarP(&conf.kwasm.assetPath, "asset-path", "a", "/assets", "Path to the asset to install")
	rootCmd.AddCommand(installCmd)
}

func runInstall(_ *cobra.Command, _ []string) error {
	// Get file or directory information.
	info, err := os.Stat(conf.kwasm.assetPath)
	if err != nil {
		return err
	}

	var files []fs.DirEntry
	// Check if the path is a directory.
	if info.IsDir() {
		files, err = os.ReadDir(conf.kwasm.assetPath)
		if err != nil {
			return err
		}
	} else {
		// If the path is not a directory, add the file to the list of files.
		files = append(files, fs.FileInfoToDirEntry(info))
		conf.kwasm.assetPath = path.Dir(conf.kwasm.assetPath)
	}

	rootFs := afero.NewOsFs()
	hostFs := afero.NewBasePathFs(rootFs, conf.host.rootPath)

	containerdConfig := containerd.NewConfig(hostFs, conf.runtime.configPath)
	shimConfig := shim.NewConfig(rootFs, hostFs, conf.kwasm.assetPath, conf.kwasm.path)

	anythingChanged := false
	for _, file := range files {
		fileName := file.Name()
		runtimeName := shim.RuntimeName(fileName)

		binPath, changed, err := shimConfig.Install(fileName)
		if err != nil {
			return fmt.Errorf("failed to install shim '%s': %w", runtimeName, err)
		}
		anythingChanged = anythingChanged || changed
		slog.Info("shim installed", "shim", runtimeName, "path", binPath, "new-version", changed)

		err = containerdConfig.AddRuntime(binPath)
		if err != nil {
			return fmt.Errorf("failed to write containerd config: %w", err)
		}
		slog.Info("shim configured", "shim", runtimeName, "path", conf.runtime.configPath)
	}

	if !anythingChanged {
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
