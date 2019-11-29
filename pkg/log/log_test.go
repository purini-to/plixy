package log

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestSetWriter(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "should have an argument logger set", args: args{logger: logger}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWriter(tt.args.logger)
			assert.Equal(t, tt.args.logger, w.logger)
		})
	}
}
