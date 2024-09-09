package godnf

import (
	"reflect"
	"testing"
)

func TestProcessOptions(t *testing.T) {
	tests := []struct {
		name string
		opt  Options
		want []string
	}{
		{
			name: "DryRun option enabled",
			opt:  Options{DryRun: true},
			want: []string{"--setopt", "tsflags=test", "--assumeyes"},
		},
		{
			name: "Verbose option enabled",
			opt:  Options{Verbose: true},
			want: []string{"--verbose", "--assumeyes"},
		},
		{
			name: "NotAssumeYes option enabled",
			opt:  Options{NotAssumeYes: true},
			want: []string{},
		},
		{
			name: "All options enabled",
			opt:  Options{DryRun: true, Verbose: true, NotAssumeYes: true},
			want: []string{"--setopt", "tsflags=test", "--verbose"},
		},
		{
			name: "No options enabled",
			opt:  Options{DryRun: false, Verbose: false, NotAssumeYes: false},
			want: []string{"--assumeyes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processOptions(&tt.opt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
