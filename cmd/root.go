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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	conf config
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "kwasm-node-installer",
	Short: "kwasm-node-installer manages containerd shims",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initializeConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&conf.runtime.name, "runtime", "r", "containerd", "Set the container runtime to configure (containerd, cri-o)")
	rootCmd.PersistentFlags().StringVarP(&conf.runtime.configPath, "runtime-config", "c", "/etc/containerd/config.toml", "Path to the runtime config file")
	rootCmd.PersistentFlags().StringVarP(&conf.kwasm.path, "kwasm-path", "k", "/opt/kwasm", "Working directory for kwasm on the host")
	rootCmd.PersistentFlags().StringVarP(&conf.host.rootPath, "host-root", "H", "/", "Path to the host root path")
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("kwasm")

	v.AddConfigPath(".")
	v.AddConfigPath("/etc")

	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	// Environment variables are prefixed with KWASM_
	v.SetEnvPrefix("KWASM")

	// - is replaced with _ in environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind env variables
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

// bindFlags binds each cobra flag to its associated viper configuration.
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
