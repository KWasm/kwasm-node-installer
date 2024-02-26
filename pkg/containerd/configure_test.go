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
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func newTestFs(fixturePath string) afero.Fs {
	baseFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join("../..", fixturePath))
	p, _ := baseFs.(*afero.BasePathFs).RealPath("/")
	fmt.Println(filepath.Abs(p))
	fs := afero.NewCopyOnWriteFs(baseFs, afero.NewMemMapFs())
	return fs
}

func TestFs(t *testing.T) {
	fs := newTestFs("testdata/containerd/valid")

	_, err := fs.Stat("/etc/containerd/config.toml")
	if err != nil {
		t.Error(err)
	}
}

func TestConfig_AddRuntime(t *testing.T) {
	type args struct {
		shimPath string
	}
	tests := []struct {
		name                     string
		args                     args
		configFile               string
		initialConfigFileContent string
		createFile               bool
		wantErr                  bool
		wantFileContent          string
	}{
		{"foobar", args{"/assets/foobar"}, "/etc/containerd/config.toml", "Hello World\n", true, false,
			`Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"
`},
		{"foobar", args{"/assets/foobar"}, "/etc/config.toml", "", false, true, ``},
		{"foobar", args{"/assets/foobar"}, "/etc/containerd/config.toml", `Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"

Foobar
`, true, false,
			`Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"

Foobar
`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if tt.createFile {
				file, err := fs.Create(tt.configFile)
				if err != nil {
					t.Fatal(err)
				}

				_, err = file.WriteString(tt.initialConfigFileContent)
				if err != nil {
					t.Fatal(err)
				}
			}

			c := &Config{
				configPath: "/etc/containerd/config.toml",
				fs:         fs,
			}
			err := c.AddRuntime(tt.args.shimPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.AddRuntime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			gotFileContent, err := afero.ReadFile(fs, tt.configFile)
			if err != nil {
				t.Fatal(err)
			}

			if string(gotFileContent) != tt.wantFileContent {
				t.Errorf("runtimeConfigFile content: %v, want %v", string(gotFileContent), tt.wantFileContent)
			}

		})
	}
}
