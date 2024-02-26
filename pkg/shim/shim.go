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

package shim

import (
	"strings"

	"github.com/spf13/afero"
)

type Config struct {
	rootFs    afero.Fs
	hostFs    afero.Fs
	assetPath string
	kwasmPath string
}

func NewConfig(rootFs afero.Fs, hostFs afero.Fs, assetPath string, kwasmPath string) *Config {
	return &Config{
		rootFs:    rootFs,
		hostFs:    hostFs,
		assetPath: assetPath,
		kwasmPath: kwasmPath,
	}
}

func RuntimeName(bin string) string {
	return strings.TrimPrefix(bin, "containerd-shim-")
}
