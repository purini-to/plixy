package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
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
			SetLogger(tt.args.logger)
			assert.Equal(t, tt.args.logger, w.logger)
		})
	}
}

func TestGetLogger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	w.logger = logger

	tests := []struct {
		name string
		want *zap.Logger
	}{
		{name: "should get a writer logger", want: logger},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLogger()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContext(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want *zap.Logger
	}{
		{
			name: "should can be get logger if logger set in the context",
			args: args{logger: logger},
			want: logger,
		},
		{
			name: "should can be get nil if logger not set in the context",
			args: args{logger: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.logger != nil {
				ctx = ToContext(ctx, tt.args.logger)
			}
			got := FromContext(ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
