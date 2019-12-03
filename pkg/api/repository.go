package api

import (
	"net/url"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
)

const (
	file = "file"
)

var repo Repository

type Repository interface {
	GetApiConfigs() ([]*config.Api, error)
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

func GetApiConfigs() ([]*config.Api, error) {
	return repo.GetApiConfigs()
}
