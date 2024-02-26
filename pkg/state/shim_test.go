package state

import (
	"testing"
)

func TestShim_MarshalJSON(t *testing.T) {
	type fields struct {
		Sha256 []byte
		Path   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"default", fields{Sha256: []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82}, Path: "/opt/kwasm/bin/containerd-shim-spin-v1"}, `{"sha256":"6da5e8f17a9bfa9cb04cf22c87b6475394ecec3af4fdc337f72d6dbf3319ea52","path":"/opt/kwasm/bin/containerd-shim-spin-v1"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shim{
				Sha256: tt.fields.Sha256,
				Path:   tt.fields.Path,
			}
			got, err := s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Shim.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("Shim.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestShim_UnmarshalJSON(t *testing.T) {
	type wants struct {
		path   string
		sha256 []byte
	}
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    wants
		wantErr bool
	}{
		{"default", args{`{"sha256":"6da5e8f17a9bfa9cb04cf22c87b6475394ecec3af4fdc337f72d6dbf3319ea52","path":"/opt/kwasm/bin/containerd-shim-spin-v1"}`}, wants{path: "/opt/kwasm/bin/containerd-shim-spin-v1", sha256: []byte{109, 165, 232, 241, 122, 155, 250, 156, 176, 76, 242, 44, 135, 182, 71, 83, 148, 236, 236, 58, 244, 253, 195, 55, 247, 45, 109, 191, 51, 25, 234, 82}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shim{}
			if err := s.UnmarshalJSON([]byte(tt.args.data)); (err != nil) != tt.wantErr {
				t.Errorf("Shim.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(s.Path) != string(tt.want.path) {
				t.Errorf("path = %v, want %v", string(s.Path), tt.want.path)
			}
			if string(s.Sha256) != string(tt.want.sha256) {
				t.Errorf("sha256 = %v, want %v", string(s.Sha256), tt.want.sha256)
			}
		})
	}
}
