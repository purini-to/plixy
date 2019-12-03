package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

type FileSystemRepository struct {
	def     *Definition
	path    string
	watcher *fsnotify.Watcher
}

func (f *FileSystemRepository) GetApiConfigs() ([]*Api, error) {
	return f.def.Apis, nil
}

func (f *FileSystemRepository) Watch(ctx context.Context, defChan chan<- *DefinitionChanged) {
	log.Debug("Start watch api definition file", zap.String("file", f.path))
	go func() {
		for {
			select {
			case event := <-f.watcher.Events:
				log.Debug("Api definition file change was detected")
				if event.Op&fsnotify.Write == fsnotify.Write {
					err := f.emitChangeApiDef(defChan)
					if err != nil {
						log.Error("Error emit write change api definition", zap.Error(err))
						continue
					}
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					if info, err := os.Stat(f.path); err != nil || info.IsDir() {
						// file not found
						continue
					}
					if err := f.watcher.Add(f.path); err != nil {
						log.Error("Could not watch file", zap.String("file", f.path))
						continue
					}
					err := f.emitChangeApiDef(defChan)
					if err != nil {
						log.Error("Error emit rename change api definition", zap.Error(err))
						continue
					}
				}
			case err := <-f.watcher.Errors:
				log.Error("Error received from file system notify", zap.Error(err))
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (f *FileSystemRepository) Close() error {
	return f.watcher.Close()
}

func (f *FileSystemRepository) validate(def *Definition) error {
	// TODO validate
	return nil
}

func (f *FileSystemRepository) parseApiDef(bytes []byte) (*Definition, error) {
	var config Definition
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal apis config file")
	}
	if err := f.validate(&config); err != nil {
		return nil, errors.Wrap(err, "invalid file system repository")
	}

	return &config, nil
}

func (f *FileSystemRepository) readApiDefFiles(path string) ([]byte, error) {
	logger := log.GetLogger().WithOptions(zap.AddStacktrace(zapcore.PanicLevel))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not read apis config file. paths: %s", path))
	}

	logger.Debug("Read apis config file", zap.String("path", path))
	return bytes, nil
}

func (f *FileSystemRepository) emitChangeApiDef(defChan chan<- *DefinitionChanged) error {
	bytes, err := f.readApiDefFiles(f.path)
	if err != nil {
		return errors.Wrap(err, "could not read the api definition file")
	}

	definition, err := f.parseApiDef(bytes)
	if err != nil {
		return errors.Wrap(err, "could not parse the api definition")
	}
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

	f := &FileSystemRepository{}
	bytes, err := f.readApiDefFiles(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error read api config file")
	}
	f.path = filePath

	definition, err := f.parseApiDef(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "error parse api config file")
	}
	f.def = definition

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a file system watcher")
	}
	f.watcher = watcher

	if err := f.watcher.Add(f.path); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not watch file. file: %s", filePath))
	}

	return f, nil
}
