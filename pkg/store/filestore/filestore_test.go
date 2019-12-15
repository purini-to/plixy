package filestore

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should be return error if file not found", func(t *testing.T) {
		got, err := New("filestore_test_not_found")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should be return error if path is directory", func(t *testing.T) {
		name, _ := ioutil.TempDir("", "filestore_test")
		defer os.Remove(name)
		got, err := New(name)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should be return error if yaml is invalid", func(t *testing.T) {
		name, _ := ioutil.TempFile("", "filestore_test")
		defer os.Remove(name.Name())

		ioutil.WriteFile(name.Name(), []byte(`
apis:
  test:invalid:test
`), 0644)

		got, err := New(name.Name())
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should be return error if definition is invalid", func(t *testing.T) {
		name, _ := ioutil.TempFile("", "filestore_test")
		defer os.Remove(name.Name())

		ioutil.WriteFile(name.Name(), []byte(`
apis:
  - name: ""
`), 0644)

		got, err := New(name.Name())
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should be return instance", func(t *testing.T) {
		name, _ := ioutil.TempFile("", "filestore_test")
		defer os.Remove(name.Name())

		ioutil.WriteFile(name.Name(), []byte(`
apis:
  - name: "test"
    proxy:
      path: "/test"
      upstream:
        target: "http://localhost:8080"
`), 0644)

		wantDef := &api.Definition{
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
		got, err := New(name.Name())
		assert.NoError(t, err)

		gotDef, _ := got.GetDefinition()
		assert.Equal(t, wantDef, gotDef)
		assert.Equal(t, name.Name(), got.filePath)
	})
}

func TestStore_Watch(t *testing.T) {
	t.Run("should be watch changed file", func(t *testing.T) {
		name, _ := ioutil.TempFile("", "filestore_test")
		defer os.Remove(name.Name())

		ioutil.WriteFile(name.Name(), []byte(`
apis:
  - name: "test"
    proxy:
      path: "/test"
      upstream:
        target: "http://localhost:8080"
`), 0644)

		wantDef := &api.Definition{
			Apis: []*api.Api{
				{
					Name: "test-update",
					Proxy: &api.Proxy{
						Path: "/test-update",
						Upstream: &api.Upstream{
							Target: "http://localhost:8080/update",
						},
					},
				},
			},
		}
		store, err := New(name.Name())
		assert.NoError(t, err)
		defer store.Close()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		defChan := make(chan *api.DefinitionChanged)
		defer close(defChan)

		err = store.Watch(ctx, 10*time.Millisecond, defChan)
		assert.NoError(t, err)

		var gotDef *api.Definition
		cnt := 0
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case def, ok := <-defChan:
					if !ok {
						return
					}
					gotDef = def.Definition
					cnt += 1
				}
			}
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			ioutil.WriteFile(name.Name(), []byte(`
apis:
  - name: "test-update"
    proxy:
      path: "/test-update"
      upstream:
        target: "http://localhost:8080/update"
`), 0644)
			time.Sleep(200 * time.Millisecond)
			cancel()
		}()
		<-ctx.Done()
		assert.Equal(t, wantDef, gotDef)
		assert.Equal(t, 1, cnt)
	})
}
