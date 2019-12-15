package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
)

const (
	file = "file"
)

var repo Repository

type Repository interface {
	GetDefinition() (*api.Definition, error)
	GetVersion() int64
}

type Watcher interface {
	Watch(ctx context.Context, defChan chan<- *api.DefinitionChanged) error
	Close() error
}

func InitRepository(dsn string) error {
	if dsn == "" {
		dsn = "file://"
	}

	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return errors.Wrap(err, "error parsing the database DSN")
	}

	switch dsnURL.Scheme {
	case file:
		log.Debug("File system based apis configuration chosen")
		repository, err := NewFileSystemRepository(dsnURL.Path)
		if err != nil {
			return errors.Wrap(err, "could not new file system repository")
		}
		repo = repository
		return nil
	default:
		return errors.New(fmt.Sprintf("The selected scheme is not supported to load api definitions. scheme: %s", dsnURL.Scheme))
	}
}

func GetDefinition() (*api.Definition, error) {
	return repo.GetDefinition()
}

func GetVersion() int64 {
	return repo.GetVersion()
}

func Watch(ctx context.Context, defChan chan<- *api.DefinitionChanged) error {
	if watcher, ok := repo.(Watcher); ok {
		return watcher.Watch(ctx, defChan)
	}
	return nil
}

func Close() error {
	if watcher, ok := repo.(Watcher); ok {
		return watcher.Close()
	}
	return nil
}
