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
	"crypto/sha256"
	"io"
	"os"
	"path"

	"github.com/kwasm/kwasm-node-installer/pkg/state"
)

func (c *Config) Install(shimName string) (string, bool, error) {
	shimPath := c.config.AssetPath(shimName)
	srcFile, err := os.OpenFile(shimPath, os.O_RDONLY, 0000)
	if err != nil {
		return "", false, err
	}
	dstFilePath := path.Join(c.config.Kwasm.Path, "bin", shimName)
	dstFilePathHost := c.config.PathWithHost(dstFilePath)

	err = os.MkdirAll(path.Dir(dstFilePathHost), 0755)
	if err != nil {
		return dstFilePath, false, err
	}

	dstFile, err := os.OpenFile(dstFilePathHost, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", false, err
	}

	st, err := state.Get(c.config)
	if err != nil {
		return "", false, err
	}
	shimSha256 := sha256.New()

	_, err = io.Copy(io.MultiWriter(dstFile, shimSha256), srcFile)
	runtimeName := RuntimeName(shimName)
	changed := st.ShimChanged(runtimeName, shimSha256.Sum(nil), dstFilePath)
	if changed {
		st.UpdateShim(runtimeName, state.Shim{
			Path:   dstFilePath,
			Sha256: shimSha256.Sum(nil),
		})
		if err := st.Write(); err != nil {
			return "", false, err
		}
	}

	return dstFilePath, changed, err
}
