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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install containerd shims",
	Run: func(cmd *cobra.Command, args []string) {

		// Get file or directory information.
		info, err := os.Stat(config.Kwasm.AssetPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		var files []fs.DirEntry
		// Check if the path is a directory.
		if info.IsDir() {
			files, err = os.ReadDir(config.Kwasm.AssetPath)
			if err != nil {
				slog.Error(err.Error())
				return
			}
		} else {
			// If the path is not a directory, add the file to the list of files.
			files = append(files, fs.FileInfoToDirEntry(info))
			config.Kwasm.AssetPath = path.Dir(config.Kwasm.AssetPath)
		}

		rootFs := afero.NewOsFs()
		hostFs := afero.NewBasePathFs(rootFs, config.Host.RootPath)

		containerdConfig := containerd.NewConfig(hostFs, config.Runtime.ConfigPath)
		shimConfig := shim.NewConfig(rootFs, hostFs, config.Kwasm.AssetPath, config.Kwasm.Path)

		anythingChanged := false
		for _, file := range files {
			fileName := file.Name()
			runtimeName := shim.RuntimeName(fileName)

			binPath, changed, err := shimConfig.Install(fileName)
			if err != nil {
				slog.Error("failed to install shim", "shim", runtimeName, "error", err)
				return
			}
			anythingChanged = anythingChanged || changed
			slog.Info("shim installed", "shim", runtimeName, "path", binPath, "new-version", changed)

			err = containerdConfig.AddRuntime(binPath)
			if err != nil {
				slog.Error("failed to write containerd config", "shim", runtimeName, "path", config.Runtime.ConfigPath, "error", err)
				return
			}
			slog.Info("shim configured", "shim", runtimeName, "path", config.Runtime.ConfigPath)
		}

		if !anythingChanged {
			slog.Info("nothing changed, nothing more to do")
			return
		}

		slog.Info("restarting containerd")
		err = containerdConfig.RestartRuntime()
		if err != nil {
			slog.Error("failed to restart containerd", "error", err)
		}
	},
}

func init() {
	installCmd.Flags().StringVarP(&config.Kwasm.AssetPath, "asset-path", "a", "/assets", "Path to the asset to install")
	rootCmd.AddCommand(installCmd)
}
