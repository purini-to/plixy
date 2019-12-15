package store

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/purini-to/plixy/pkg/store/filestore"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/api"
	"github.com/purini-to/plixy/pkg/log"
)

const defaultDSN = "file://"

const (
	fileSchema = "file"
)

// Store defines the behavior of a proxy specs.
type Store interface {
	GetDefinition() (*api.Definition, error)
	Watch(ctx context.Context, interval time.Duration, defChan chan<- *api.DefinitionChanged) error
	Close() error
}

// Build creates a store instance that will depend on your given DSN.
func Build(dsn string) (Store, error) {
	if dsn == "" {
		dsn = defaultDSN
	}

	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing the store DSN")
	}

	var store Store
	switch dsnURL.Scheme {
	case fileSchema:
		log.Debug("File system based apis configuration chosen")
		path := dsnURL.Path
		if path == "" {
			path = "./plixy.yaml"
		}
		store, err = filestore.New(path)
		if err != nil {
			return nil, errors.Wrap(err, "could not new file system store")
		}
	default:
		return nil, errors.New(fmt.Sprintf("The selected scheme is not supported to load api definitions. scheme: %s", dsnURL.Scheme))
	}

	return store, nil
}
