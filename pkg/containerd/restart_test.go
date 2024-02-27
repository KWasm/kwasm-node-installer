package containerd

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-ps"
)

type mockProcess struct {
	executable string
	pid        int
	ppid       int
}

func (p *mockProcess) Executable() string {
	return p.executable
}

func (p *mockProcess) Pid() int {
	return p.pid
}

func (p *mockProcess) PPid() int {
	return p.ppid
}

func Test_getPid(t *testing.T) {
	tests := []struct {
		name             string
		psProccessesMock func() ([]ps.Process, error)
		want             int
		wantErr          bool
	}{
		{"no containerd process found", func() ([]ps.Process, error) {
			return []ps.Process{}, nil
		}, -1, true},
		{"single containerd process found", func() ([]ps.Process, error) {
			return []ps.Process{
				&mockProcess{executable: "containerd", pid: 123},
			}, nil
		}, 123, false},
		{"multiple containerd processes found", func() ([]ps.Process, error) {
			return []ps.Process{
				&mockProcess{executable: "containerd"},
				&mockProcess{executable: "containerd"},
			}, nil
		}, 0, true},
		{"error getting processes", func() ([]ps.Process, error) {
			return nil, fmt.Errorf("error getting processes")
		}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						// If the test case does not expect an error, re-panic
						panic(r)
					}
				}
			}()

			psProcesses = tt.psProccessesMock
			got, err := getPid()
			if (err != nil) != tt.wantErr {
				t.Errorf("getPid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getPid() = %v, want %v", got, tt.want)
			}
		})
	}
}
