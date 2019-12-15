package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/purini-to/plixy/pkg/plugin"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"
)

type FileSystemRepository struct {
	sync.RWMutex
	def     *api.Definition
	path    string
	ticker  *time.Ticker
	version int64
}

func (f *FileSystemRepository) GetDefinition() (*api.Definition, error) {
	return f.def, nil
}

func (f *FileSystemRepository) GetVersion() int64 {
	return f.version
}

func (f *FileSystemRepository) Watch(ctx context.Context, defChan chan<- *api.DefinitionChanged) error {
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

				f.checkNewDefVersion(info, defChan)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (f *FileSystemRepository) checkNewDefVersion(info os.FileInfo, defChan chan<- *api.DefinitionChanged) {
	f.Lock()
	defer f.Unlock()

	if f.version >= info.ModTime().Unix() {
		return
	}
	f.version = info.ModTime().Unix()
	log.Info("Api definition file change detected")

	err := f.emitChangeApiDef(defChan)
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

func (f *FileSystemRepository) validate(def *api.Definition) error {
	isValid, err := def.Validate()
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("invalid api definition")
	}

	plugins := make([]*api.Plugin, 0)
	for _, a := range def.Apis {
		for _, p := range a.Plugins {
			plugins = append(plugins, p)
		}
	}
	if err := plugin.ValidateConfig(plugins); err != nil {
		return err
	}
	return nil
}

func (f *FileSystemRepository) parseApiDef(bytes []byte) (*api.Definition, error) {
	var definition api.Definition
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

func (f *FileSystemRepository) emitChangeApiDef(defChan chan<- *api.DefinitionChanged) error {
	bytes, err := f.readApiDefFiles(f.path)
	if err != nil {
		return errors.Wrap(err, "could not read the api definition file")
	}

	definition, err := f.parseApiDef(bytes)
	if err != nil {
		return errors.Wrap(err, "could not parse the api definition")
	}
	f.def = definition

	defChan <- &api.DefinitionChanged{
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
	f.version = info.ModTime().Unix()
	f.def = definition

	return f, nil
}
