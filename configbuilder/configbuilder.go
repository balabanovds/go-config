package configbuilder

import (
	"errors"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var (
	ErrFilenameEmpty   = errors.New("empty file name")
	ErrFileWrongFormat = errors.New("file wrong format")
)

type ConfigBuilder struct {
	k      *koanf.Koanf
	errors []error
}

func New() *ConfigBuilder {
	return &ConfigBuilder{k: koanf.New(".")}
}

func (c *ConfigBuilder) LoadToml(filename string) *ConfigBuilder {
	return c.loadFile(filename, toml.Parser())
}

func (c *ConfigBuilder) LoadJSON(filename string) *ConfigBuilder {
	return c.loadFile(filename, json.Parser())
}

func (c *ConfigBuilder) LoadEnv(prefix, keyDelimiter, valueDelimiter string) *ConfigBuilder {
	cbFunc := func(key string, value string) (string, interface{}) {
		key = strings.Replace(strings.ToLower(strings.TrimPrefix(key, prefix+keyDelimiter)), "_", ".", -1)

		if valueDelimiter == "" {
			return key, value
		}

		value = strings.Trim(value, valueDelimiter)
		values := strings.Split(value, valueDelimiter)

		if len(values) == 1 {
			return key, values[0]
		}

		return key, values
	}

	err := c.k.Load(env.ProviderWithValue(prefix, ".", cbFunc), nil)
	if err != nil {
		c.addError(err)
	}

	return c
}

func (c *ConfigBuilder) ToStruct(cfg interface{}) error {
	if len(c.errors) != 0 {
		err := c.errors[0]
		c.errors = []error{}
		return err
	}

	if err := c.k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}

func (c *ConfigBuilder) loadFile(filename string, parser koanf.Parser) *ConfigBuilder {
	if filename == "" {
		c.addError(ErrFilenameEmpty)
		return c
	}

	fi, err := os.Stat(filename)
	if err != nil {
		c.addError(err)
		return c
	}

	if !fi.Mode().IsRegular() {
		c.addError(ErrFileWrongFormat)
		return c
	}

	if err := c.k.Load(file.Provider(filename), parser); err != nil {
		c.addError(err)
	}

	return c
}

func (c *ConfigBuilder) addError(err error) {
	c.errors = append(c.errors, err)
}
