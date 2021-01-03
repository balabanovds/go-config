package configbuilder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	configbuilder "github.com/balabanovds/go-config"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

var (
	intVar = 1
	strVar = "var2"
	arrVar = [...]int{1, 2, 3}

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
	tomlFile string
	jsonFile string
}

func (s *TestSuite) SetupTest() {
	tomlConfig := fmt.Sprintf(`[toml]
int = %d
str = %s
arr = [1,2,3]
`, intVar, strVar)

	jsonConfig := fmt.Sprintf(`{
	"json": {
		"int": %d,
		"str": %s,
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
	cb := configbuilder.New()

	cfg := new(Cfg)

	err := cb.LoadToml(s.tomlFile).ToStruct(&cfg)
	require.NoError(s.T(), err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
