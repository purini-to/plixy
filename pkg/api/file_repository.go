package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"
)

type FileSystemRepository struct {
	sync.RWMutex
	def    *Definition
	path   string
	ticker *time.Ticker
}

func (f *FileSystemRepository) GetDefinition() (*Definition, error) {
	return f.def, nil
}

func (f *FileSystemRepository) Watch(ctx context.Context, defChan chan<- *DefinitionChanged) error {
	f.ticker = time.NewTicker(config.Global.WatchInterval)

	log.Debug("Start watch api definition file", zap.String("file", f.path))
	go func() {
		for {
			select {
			case <-f.ticker.C:
				info, err := os.Stat(f.path)
				if err != nil || info.IsDir() {
					// file not found
					continue
				}

				f.Lock()
				f.checkNewDefVersion(info, defChan)
				f.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (f *FileSystemRepository) checkNewDefVersion(info os.FileInfo, defChan chan<- *DefinitionChanged) {
	if f.def.Version >= info.ModTime().Unix() {
		return
	}
	log.Info("Api definition file change detected")

	err := f.emitChangeApiDef(defChan, info.ModTime().Unix())
	if err != nil {
		log.Error("Error emit rename change api definition", zap.Error(err))
		return
	}
}

func (f *FileSystemRepository) Close() error {
	if f.ticker == nil {
		return nil
	}
	f.ticker.Stop()
	return nil
}

func (f *FileSystemRepository) validate(def *Definition) error {
	// TODO validate
	return nil
}

func (f *FileSystemRepository) parseApiDef(bytes []byte) (*Definition, error) {
	var definition Definition
	if err := yaml.Unmarshal(bytes, &definition); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal apis definition file")
	}
	if err := f.validate(&definition); err != nil {
		return nil, errors.Wrap(err, "invalid file system repository")
	}

	return &definition, nil
}

func (f *FileSystemRepository) readApiDefFiles(path string) ([]byte, error) {
	logger := log.GetLogger().WithOptions(zap.AddStacktrace(zapcore.PanicLevel))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not read apis definition file. paths: %s", path))
	}

	logger.Debug("Read apis definition file", zap.String("path", path))
	return bytes, nil
}

func (f *FileSystemRepository) emitChangeApiDef(defChan chan<- *DefinitionChanged, version int64) error {
	bytes, err := f.readApiDefFiles(f.path)
	if err != nil {
		return errors.Wrap(err, "could not read the api definition file")
	}

	definition, err := f.parseApiDef(bytes)
	if err != nil {
		return errors.Wrap(err, "could not parse the api definition")
	}
	definition.Version = version
	f.def = definition

	defChan <- &DefinitionChanged{
		Definition: definition,
	}
	return nil
}

func NewFileSystemRepository(filePath string) (*FileSystemRepository, error) {
	if filePath == "" {
		filePath = "./plixy.yaml"
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return nil, errors.Wrap(err, fmt.Sprintf("not found api definition file. paths: %s", filePath))
	}

	f := &FileSystemRepository{}
	bytes, err := f.readApiDefFiles(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error read api definition file")
	}
	f.path = filePath

	definition, err := f.parseApiDef(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "error parse api definition file")
	}
	definition.Version = info.ModTime().Unix()
	f.def = definition

	return f, nil
}
