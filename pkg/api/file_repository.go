package api

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"
)

type FileSystemRepository struct {
	def *config.ApiDefinition
}

func (f *FileSystemRepository) GetApiConfigs() ([]*config.Api, error) {
	return f.def.Apis, nil
}

func (f *FileSystemRepository) validate() error {
	// TODO validate
	return nil
}

func NewFileSystemRepository(filePath string) (*FileSystemRepository, error) {
	def, err := loadApiDef(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a file system repository")
	}

	f := &FileSystemRepository{def: def}
	if err := f.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid file system repository")
	}

	return f, nil
}

func loadApiDef(filePath string) (*config.ApiDefinition, error) {
	paths := []string{filePath}
	if paths[0] == "" {
		paths = []string{"./plixy.yaml", "/etc/plixy/plixy.yaml"}
	}

	bytes, err := readApiDefFiles(paths)
	if err != nil {
		return nil, errors.Wrap(err, "failed read apis config file")
	}

	var config config.ApiDefinition
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal apis config file")
	}

	return &config, nil
}

func readApiDefFiles(paths []string) ([]byte, error) {
	logger := log.GetLogger().WithOptions(zap.AddStacktrace(zapcore.PanicLevel))
	for _, p := range paths {
		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			logger.Warn("No apis config file found", zap.String("path", p), zap.Error(err))
			continue
		}

		logger.Debug("Read apis config file", zap.String("path", p))
		return bytes, nil
	}

	return nil, errors.New(fmt.Sprintf("could not read apis config file. paths: %s", paths))
}
