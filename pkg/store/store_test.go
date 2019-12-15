package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/purini-to/plixy/pkg/store/filestore"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Run("should be building file store if schema is file", func(t *testing.T) {
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

		got, err := Build("file://" + name.Name())
		assert.NoError(t, err)
		assert.IsType(t, &filestore.Store{}, got)
	})

	t.Run("should be return error if schema is unknown", func(t *testing.T) {
		got, err := Build("unknown://test")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should be return error if dsn is invalid url", func(t *testing.T) {
		got, err := Build(":")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
