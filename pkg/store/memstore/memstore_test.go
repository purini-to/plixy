package memstore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/purini-to/plixy/pkg/api"
)

func TestStore_Close(t *testing.T) {
	t.Run("should be no error in calling close()", func(t *testing.T) {
		s := New()
		err := s.Close()
		assert.NoError(t, err)
	})
}

func TestStore_GetDefinition(t *testing.T) {
	t.Run("should be get api.Definition", func(t *testing.T) {
		want := &api.Definition{}
		s := New()
		s.def = want
		got, err := s.GetDefinition()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestStore_SetDefinition(t *testing.T) {
	t.Run("should be set api.Definition in memstore", func(t *testing.T) {
		want := &api.Definition{
			Apis: []*api.Api{
				{
					Name: "test",
					Proxy: &api.Proxy{
						Path: "/test",
						Upstream: &api.Upstream{
							Target: "http://localhost:8080",
						},
					},
				},
			},
		}
		s := New()
		err := s.SetDefinition(want)
		assert.NoError(t, err)
		assert.Equal(t, want, s.def)
	})

	t.Run("should be no set api.Definition if definition has error", func(t *testing.T) {
		s := New()
		err := s.SetDefinition(&api.Definition{})
		assert.Error(t, err)
		assert.Nil(t, s.def)
	})
}
