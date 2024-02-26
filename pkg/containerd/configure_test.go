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

func newFixtureFs(fixturePath string) afero.Fs {
	baseFs := afero.NewBasePathFs(afero.NewOsFs(), filepath.Join("../..", fixturePath))
	p, _ := baseFs.(*afero.BasePathFs).RealPath("/")
	fmt.Println(filepath.Abs(p))
	fs := afero.NewCopyOnWriteFs(baseFs, afero.NewMemMapFs())
	return fs
}

func TestConfig_AddRuntime(t *testing.T) {
	type fields struct {
		hostFs     afero.Fs
		configPath string
	}
	type args struct {
		shimPath string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantFileErr     bool
		wantFileContent string
	}{
		{"missing shim config", fields{
			hostFs:     newFixtureFs("testdata/containerd/missing-containerd-shim-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, false, false, `[plugins]
  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false
  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]
  [plugins."io.containerd.gc.v1.scheduler"]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = 0
    startup_delay = "100ms"
  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = true
  [plugins."io.containerd.service.v1.tasks-service"]
    blockio_config_file = ""
    rdt_config_file = ""

# KWASM runtime config for spin-v1
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin-v1]
runtime_type = "/opt/kwasm/bin/containerd-shim-spin-v1"
`},
		{"missing config", fields{
			hostFs:     newFixtureFs("testdata/containerd/missing-containerd-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, true, true, ``},
		{"existing shim config", fields{
			hostFs:     newFixtureFs("testdata/containerd/existing-containerd-shim-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, false, false, `[plugins]
  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false
  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]
  [plugins."io.containerd.gc.v1.scheduler"]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = 0
    startup_delay = "100ms"
  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = true
  [plugins."io.containerd.service.v1.tasks-service"]
    blockio_config_file = ""
    rdt_config_file = ""

# KWASM runtime config for spin-v1
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin-v1]
runtime_type = "/opt/kwasm/bin/containerd-shim-spin-v1"
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				hostFs:     tt.fields.hostFs,
				configPath: tt.fields.configPath,
			}
			if err := c.AddRuntime(tt.args.shimPath); (err != nil) != tt.wantErr {
				t.Errorf("Config.AddRuntime() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotContent, err := afero.ReadFile(c.hostFs, c.configPath)
			if (err != nil) != tt.wantFileErr {
				t.Errorf("read %s error = %v, wantFileErr %v", c.configPath, err, tt.wantFileErr)
				return
			}

			if string(gotContent) != tt.wantFileContent {
				t.Errorf("file content %s got = %s, want = %s", c.configPath, string(gotContent), tt.wantFileContent)
			}
		})
	}
}

func TestConfig_RemoveRuntime(t *testing.T) {
	type fields struct {
		hostFs     afero.Fs
		configPath string
	}
	type args struct {
		shimPath string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantFileErr     bool
		wantFileContent string
	}{
		{"missing shim config", fields{
			hostFs:     newFixtureFs("testdata/containerd/missing-containerd-shim-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, false, false, `[plugins]
  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false
  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]
  [plugins."io.containerd.gc.v1.scheduler"]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = 0
    startup_delay = "100ms"
  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = true
  [plugins."io.containerd.service.v1.tasks-service"]
    blockio_config_file = ""
    rdt_config_file = ""
`},
		{"missing config", fields{
			hostFs:     newFixtureFs("testdata/containerd/missing-containerd-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, true, true, ``},
		{"existing shim config", fields{
			hostFs:     newFixtureFs("testdata/containerd/existing-containerd-shim-config"),
			configPath: "/etc/containerd/config.toml",
		}, args{"/opt/kwasm/bin/containerd-shim-spin-v1"}, false, false, `[plugins]
  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false
  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]
  [plugins."io.containerd.gc.v1.scheduler"]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = 0
    startup_delay = "100ms"
  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = true
  [plugins."io.containerd.service.v1.tasks-service"]
    blockio_config_file = ""
    rdt_config_file = ""
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				hostFs:     tt.fields.hostFs,
				configPath: tt.fields.configPath,
			}
			if err := c.RemoveRuntime(tt.args.shimPath); (err != nil) != tt.wantErr {
				t.Errorf("Config.RemoveRuntime() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotContent, err := afero.ReadFile(c.hostFs, c.configPath)
			if (err != nil) != tt.wantFileErr {
				t.Errorf("read %s error = %v, wantFileErr %v", c.configPath, err, tt.wantFileErr)
				return
			}

			if string(gotContent) != tt.wantFileContent {
				t.Errorf("file content %s got = %s, want = %s", c.configPath, string(gotContent), tt.wantFileContent)
			}
		})
	}
}