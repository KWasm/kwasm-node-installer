package state

import (
	"fmt"
	"path/filepath"
	"reflect"
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

func TestGet(t *testing.T) {
	type args struct {
		fs        afero.Fs
		kwasmPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *state
		wantErr bool
	}{
		{"existing state", args{newFixtureFs("testdata/containerd/existing-containerd-shim-config"), "/opt/kwasm"}, &state{Shims: map[string]*Shim{
			"spin-v1": {Sha256: []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82}, Path: "/opt/kwasm/bin/containerd-shim-spin-v1"},
		}}, false},
		{"missing state", args{newFixtureFs("testdata/containerd/missing-containerd-shim-config"), "/opt/kwasm"}, &state{Shims: map[string]*Shim{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.fs, tt.args.kwasmPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Shims, tt.want.Shims) {
				t.Errorf("Get() = %v, want %v", got.Shims, tt.want.Shims)
			}
		})
	}
}

func TestShimChanged(t *testing.T) {
	state := &state{
		Shims: map[string]*Shim{
			"spin-v1": {
				Sha256: []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82},
				Path:   "/opt/kwasm/bin/containerd-shim-spin-v1",
			},
		},
	}

	tests := []struct {
		name       string
		shimName   string
		sha256     []byte
		path       string
		wantResult bool
	}{
		{
			name:       "existing shim, same sha256 and path",
			shimName:   "spin-v1",
			sha256:     []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82},
			path:       "/opt/kwasm/bin/containerd-shim-spin-v1",
			wantResult: false,
		},
		{
			name:       "existing shim, different sha256",
			shimName:   "spin-v1",
			sha256:     []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 83},
			path:       "/opt/kwasm/bin/containerd-shim-spin-v1",
			wantResult: true,
		},
		{
			name:       "existing shim, different path",
			shimName:   "spin-v1",
			sha256:     []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82},
			path:       "/opt/kwasm/bin/containerd-shim-spin-v2",
			wantResult: true,
		},
		{
			name:       "non-existing shim",
			shimName:   "non-existing",
			sha256:     []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82},
			path:       "/opt/kwasm/bin/containerd-shim-spin-v1",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := state.ShimChanged(tt.shimName, tt.sha256, tt.path)
			if result != tt.wantResult {
				t.Errorf("ShimChanged() = %v, want %v", result, tt.wantResult)
			}
		})
	}
}
