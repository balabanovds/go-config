package configbuilder

import (
	"errors"
	"os"
	"strings"

	"github.com/knadh/koanf/providers/env"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

var (
	ErrFilenameEmpty   = errors.New("empty file name")
	ErrFileWrongFormat = errors.New("file wrong format")
	ErrBuilder         = errors.New("builder contains error")
)

type config struct {
	k      *koanf.Koanf
	errors []error
}

func New() *config {
	return &config{k: koanf.New(".")}
}

func (c *config) LoadToml(filename string) *config {
	return c.loadFile(filename, toml.Parser())
}

func (c *config) LoadJSON(filename string) *config {
	return c.loadFile(filename, json.Parser())
}

func (c *config) LoadEnv(prefix string) *config {
	err := c.k.Load(env.Provider(prefix, "_", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil)
	if err != nil {
		c.addError(err)
	}

	return c
}

func (c *config) ToStruct(cfg interface{}) error {
	if len(c.errors) != 0 {
		return ErrBuilder
	}

	if err := c.k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}

func (c *config) loadFile(filename string, parser koanf.Parser) *config {
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

func (c *config) addError(err error) {
	c.errors = append(c.errors, err)
}
