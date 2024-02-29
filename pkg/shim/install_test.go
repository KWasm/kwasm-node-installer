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

func TestConfig_Install(t *testing.T) {
	type fields struct {
		rootFs    afero.Fs
		hostFs    afero.Fs
		assetPath string
		kwasmPath string
	}
	type args struct {
		shimName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			"no changes to shim",
			fields{
				newFixtureFs("testdata"),
				newFixtureFs("testdata/shim"),
				"/assets",
				"/opt/kwasm"},
			args{"containerd-shim-spin-v1"},
			"/opt/kwasm/bin/containerd-shim-spin-v1",
			false,
			false,
		},
		{
			"install new shim over old",
			fields{
				newFixtureFs("testdata"),
				newFixtureFs("testdata/shim"),
				"/assets",
				"/opt/kwasm"},
			args{"containerd-shim-slight-v1"},
			"/opt/kwasm/bin/containerd-shim-slight-v1",
			true,
			false,
		},
		{
			"unable to find new shim",
			fields{
				afero.NewMemMapFs(),
				newFixtureFs("testdata/shim"),
				"/assets",
				"/opt/kwasm"},
			args{"some-shim"},
			"",
			false,
			true,
		},
		{
			"unable to write to hostFs",
			fields{
				newFixtureFs("testdata"),
				afero.NewReadOnlyFs(newFixtureFs("testdata/shim")),
				"/assets",
				"/opt/kwasm"},
			args{"containerd-shim-spin-v1"},
			"/opt/kwasm/bin/containerd-shim-spin-v1",
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				rootFs:    tt.fields.rootFs,
				hostFs:    tt.fields.hostFs,
				assetPath: tt.fields.assetPath,
				kwasmPath: tt.fields.kwasmPath,
			}

			got, got1, err := c.Install(tt.args.shimName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Config.Install() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Config.Install() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
