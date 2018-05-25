package _internal

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	// "path/filepath"
)

type Config struct {
	ExcludedRegexps []string `yaml:"excluded_regexps"` // 需要追加监听的文件后缀名字，默认是'.go'，
	ExternPaths     []string `yaml:"extern_paths"`     // 额外需要监听的路径
	ExcludedPaths   []string `yaml:"excluded_paths"`   // 不需要监听的目录
	Envs            []string `yaml:"envs"`             // 执行时追加的环境变量
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"gw config: \n\texcluded_regexps: %s\n\textern_paths: %s\n\texcluded_paths: %s\n\tenvs: %s\n",
		c.ExcludedRegexps,
		c.ExternPaths,
		c.ExcludedPaths,
		c.Envs,
	)
}

var (
	cfg         Config
	defaultConf *Config
)

func GetInstance() *Config {
	return &cfg
}

func ParseConfig(fname string) {
	var (
		yamlFile []byte
		err      error
	)
	if yamlFile, err = ioutil.ReadFile(fname); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		panic(err)
	}
}

func OutputDefaultConf(outpath string) error {
	defaultConf = &Config{
		ExcludedRegexps: []string{
			".gitignore$",
			".yml$",
			".txt$",
		},
		ExternPaths: []string{},
		Envs: []string{
			"GOROOT=/path/to/your/goroot",
			"GOPATH=/path/to/your/gopath",
		},
		ExcludedPaths: []string{
			"vendor",
		},
	}

	bs, _ := yaml.Marshal(defaultConf)

	if err := ioutil.WriteFile(outpath, bs, 0644); err != nil {
		return err
	}
	return nil
}
