package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yeqown/go-watcher/internal"
	"github.com/yeqown/go-watcher/utils"

	"github.com/silenceper/log"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v1"
)

var (
	exit           chan bool
	cmdArgs, paths []string
	cfg            *Config
	watcher        *internal.Watcher
)

func init() {
	exit = make(chan bool, 10)
}

func initCommand() *cli.App {
	app := cli.NewApp()

	app.Name = "go-watcher"
	app.Version = "1.1.0"
	app.Author = "yeqown@gmail.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./config.yml",
			Usage: "load configuration from `FILE`, --default=./config.yml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "generate a config.yml to specified postion",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "./config.yml",
					Usage: "set output file name and position",
				},
			},
			Action: func(c *cli.Context) error {
				return generateDefaultConfigFile(c.String("output"))
			},
		},
		{
			Name:  "run",
			Usage: "execute a command, and watch the files, if any change to these files, the command will reload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "command, e",
					Value: "echo 'go-watcher'",
					Usage: "set command to execute for hot-reloading",
				},
			},
			Action: func(c *cli.Context) error {
				var (
					cfg *Config
					err error
				)
				if cfg, err = loadConfigFile(c.GlobalString("config")); err != nil {
					return err
				}
				pwd, _ := os.Getwd()

				// passing config
				paths = append(paths, cfg.ExternPaths...)
				utils.WalkDirectory(pwd, cfg.ExcludedPaths, &paths, true)
				if watcher, err = internal.
					NewWatcher(paths, exit, []string{"go"}, cfg.ExcludedRegexps); err != nil {
					return err
				}
				parsed := strings.Split(c.String("command"), " ")
				watcher.SetCommand(parsed[0], parsed[1:], cfg.Envs)
				go watcher.Watching() // start watching

				return nil
			},
		},
	}

	return app
}

// TODO: handle command and main goroutine
func main() {

	app := initCommand()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}

	// handle os signal
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT:
				close(exit)
			case syscall.SIGQUIT:
				close(exit)
			case syscall.SIGHUP:
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// main concurency keep going
	for {
		select {
		case <-exit:
			log.Info("go-watcher Exit!")
			watcher.Exit()
			os.Exit(2)
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

// Config ...
type Config struct {
	ExcludedRegexps []string `yaml:"excluded_regexps"` // 需要追加监听的文件后缀名字，默认是'.go'，
	ExternPaths     []string `yaml:"extern_paths"`     // 额外需要监听的路径
	ExcludedPaths   []string `yaml:"excluded_paths"`   // 不需要监听的目录
	Envs            []string `yaml:"envs"`             // 执行时追加的环境变量
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"go-watcher config: \n\texcluded_regexps: %s\n\textern_paths: %s\n\texcluded_paths: %s\n\tenvs: %s\n",
		c.ExcludedRegexps,
		c.ExternPaths,
		c.ExcludedPaths,
		c.Envs,
	)
}

// loadConfigFile ...
func loadConfigFile(filename string) (cfg *Config, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filename); err == nil {
		if err = yaml.Unmarshal(b, cfg); err == nil {
			return cfg, nil
		}
	}
	return nil, err
}

// generateDefaultConfigFile ...
func generateDefaultConfigFile(outpath string) error {
	c := &Config{
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

	bs, _ := yaml.Marshal(c)
	if err := ioutil.WriteFile(outpath, bs, 0644); err != nil {
		return err
	}
	return nil
}
