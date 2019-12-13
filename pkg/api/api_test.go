package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	yaml "gopkg.in/yaml.v2"
)

func TestUpstream_UnmarshalYAML(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    *Upstream
		wantErr bool
	}{
		{
			name: "should be set a value in the target if yaml has a target on the key",
			args: args{in: `
target: "http://localhost:9002/apis/v1"
test: 1
`},
			want:    &Upstream{Target: "http://localhost:9002/apis/v1"},
			wantErr: false,
		},
		{
			name: "should be set vars if target has a path variables",
			args: args{in: `
target: "http://localhost:9002/apis/v1/users/{userId}/tasks/{taskId}"
`},
			want: &Upstream{
				Target: "http://localhost:9002/apis/v1/users/{userId}/tasks/{taskId}",
				Vars:   []string{"userId", "taskId"},
			},
			wantErr: false,
		},
		{
			name: "should be error if yaml invalid",
			args: args{in: `
target: 'test""
`},
			want:    &Upstream{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Upstream{}
			err := yaml.Unmarshal([]byte(tt.args.in), u)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, u)
			}
		})
	}
}
