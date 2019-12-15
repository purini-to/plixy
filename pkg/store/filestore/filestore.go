package filestore

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"

	"github.com/purini-to/plixy/pkg/api"
	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/purini-to/plixy/pkg/store/memstore"
)

// Store represents a file system store.
type Store struct {
	sync.RWMutex
	*memstore.Store
	filePath string
	ticker   *time.Ticker
	version  int64
}

// New creates a file system store.
// return error if not found file or file is directory.
func New(filePath string) (*Store, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "not found file")
	}
	if info.IsDir() {
		return nil, errors.New("path is directory")
	}

	s := &Store{
		Store:    memstore.New(),
		filePath: filePath,
	}

	def, err := s.readApiDefFile()
	if err != nil {
		return nil, err
	}
	if err = s.SetDefinition(def); err != nil {
		return nil, err
	}

	s.version = info.ModTime().UnixNano()
	return s, nil
}

// Watch watches for changes on the file.
func (s *Store) Watch(ctx context.Context, interval time.Duration, defChan chan<- *api.DefinitionChanged) error {
	s.ticker = time.NewTicker(interval)

	log.Debug("Start watch api definition file", zap.String("filePath", s.filePath))
	go func() {
		for {
			select {
			case <-s.ticker.C:
				info, err := os.Stat(s.filePath)
				if err != nil || info.IsDir() {
					// file not found
					continue
				}

				s.emitDefinitionChanged(info.ModTime().UnixNano(), defChan)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Close terminates watch the session.
func (s *Store) Close() error {
	if s.ticker == nil {
		return s.Store.Close()
	}
	s.ticker.Stop()
	return s.Store.Close()
}

func (s *Store) readApiDefFile() (*api.Definition, error) {
	bytes, err := s.readFile()
	if err != nil {
		return nil, err
	}
	return s.parseApiDef(bytes)
}

func (s *Store) readFile() ([]byte, error) {
	bytes, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	return bytes, nil
}

func (s *Store) parseApiDef(bytes []byte) (*api.Definition, error) {
	var def api.Definition
	if err := yaml.Unmarshal(bytes, &def); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal apis definition")
	}
	return &def, nil
}

func (s *Store) emitDefinitionChanged(version int64, defChan chan<- *api.DefinitionChanged) {
	s.Lock()
	defer s.Unlock()

	if s.version >= version {
		return
	}
	log.Info(fmt.Sprintf("Api definition file change detected. %d => %d", s.version, version))
	s.version = version

	def, err := s.readApiDefFile()
	if err != nil {
		log.Error("Could not read api definition file", zap.Error(err))
		return
	}

	if err = s.SetDefinition(def); err != nil {
		log.Error("Failed set api definition", zap.Error(err))
		return
	}

	defChan <- &api.DefinitionChanged{
		Definition: def,
	}
}
