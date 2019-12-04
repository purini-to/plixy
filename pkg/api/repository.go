package api

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
)

const (
	file = "file"
)

var repo Repository

type Repository interface {
	GetApiConfigs() ([]*Api, error)
}

type Watcher interface {
	Watch(ctx context.Context, defChan chan<- *DefinitionChanged) error
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
		return errors.New("The selected scheme is not supported to load api definitions")
	}
}

func GetApiConfigs() ([]*Api, error) {
	return repo.GetApiConfigs()
}

func Watch(ctx context.Context, defChan chan<- *DefinitionChanged) error {
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
