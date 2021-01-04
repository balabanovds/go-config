package configbuilder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/balabanovds/goutils/configbuilder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	intVar = 1
	strVar = "var2"
	arrVar = []int{1, 2, 3}

	prefix = "TEST"
)

type Cfg struct {
	Toml InnerCfg `koanf:"toml"`
	Json InnerCfg `koanf:"json"`
	Env  InnerCfg `koanf:"env"`
}

type InnerCfg struct {
	Int int    `koanf:"int"`
	Str string `koanf:"str"`
	Arr []int  `koanf:"arr"`
}

type TestSuite struct {
	suite.Suite
	cfg      *Cfg
	cb       *configbuilder.ConfigBuilder
	tomlFile string
	jsonFile string
}

func (s *TestSuite) SetupTest() {
	tomlConfig := fmt.Sprintf(`[toml]
int = %d
str = "%s"
arr = [1,2,3]
`, intVar, strVar)

	jsonConfig := fmt.Sprintf(`{
	"json": {
		"int": %d,
		"str": "%s",
		"arr": [1,2,3]
	}
}`, intVar, strVar)

	f, err := ioutil.TempFile("/tmp", "toml")
	require.NoError(s.T(), err)

	s.tomlFile = f.Name()
	_, err = f.WriteString(tomlConfig)
	require.NoError(s.T(), err)

	f, err = ioutil.TempFile("/tmp", "json")
	require.NoError(s.T(), err)

	s.jsonFile = f.Name()
	_, err = f.WriteString(jsonConfig)
	require.NoError(s.T(), err)

	s.setEnv("int", "1")
	s.setEnv("str", "var2")
	s.setEnv("arr", "1,2,3")

	s.cb = configbuilder.New()
	s.cfg = new(Cfg)
}

func (s *TestSuite) setEnv(key, value string) {
	err := os.Setenv(fmt.Sprintf("%s_ENV_%s", prefix, key), value)
	require.NoError(s.T(), err)
}

func (s *TestSuite) TearDownTest() {
	assert.NoError(s.T(), os.Remove(s.tomlFile))
	assert.NoError(s.T(), os.Remove(s.jsonFile))
}

func (s *TestSuite) TestToml() {
	want := &Cfg{
		Toml: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
	}

	err := s.cb.LoadToml(s.tomlFile).ToStruct(s.cfg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, s.cfg)
}

func (s *TestSuite) TestJson() {
	want := &Cfg{
		Json: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
	}

	err := s.cb.LoadJSON(s.jsonFile).ToStruct(s.cfg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, s.cfg)
}

func (s *TestSuite) TestEnv() {
	want := &Cfg{
		Env: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
	}

	err := s.cb.LoadEnv(prefix, "_", ",").ToStruct(s.cfg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, s.cfg)
}

func (s *TestSuite) TestAll() {
	want := &Cfg{
		Toml: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
		Json: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
		Env: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
	}

	err := s.cb.LoadToml(s.tomlFile).
		LoadJSON(s.jsonFile).
		LoadEnv(prefix, "_", ",").
		ToStruct(s.cfg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, s.cfg)
}

func (s *TestSuite) TestEnvArrEndsAndStartsWithComma() {
	s.setEnv("arr", ",,1,2,3,,,")
	want := &Cfg{
		Env: InnerCfg{
			Int: intVar,
			Str: strVar,
			Arr: arrVar,
		},
	}

	err := s.cb.LoadEnv(prefix, "_", ",").ToStruct(s.cfg)
	require.NoError(s.T(), err)
	require.Equal(s.T(), want, s.cfg)
}

func (s *TestSuite) TestErrors() {
	s.Run("empty filename", func() {
		err := s.cb.LoadToml("").ToStruct(s.cfg)
		require.Error(s.T(), err)
		require.EqualError(s.T(), err, configbuilder.ErrFilenameEmpty.Error())
	})

	s.Run("dir as filename", func() {
		err := s.cb.LoadJSON(".").ToStruct(s.cfg)
		require.Error(s.T(), err)
		require.EqualError(s.T(), err, configbuilder.ErrFileWrongFormat.Error())
	})

	s.Run("no loaders", func() {
		err := s.cb.ToStruct(s.cfg)
		require.Error(s.T(), err)
		require.EqualError(s.T(), err, configbuilder.ErrNoLoader.Error())
	})
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
