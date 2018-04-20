/*
 * entry
 */
package main

import (
	"errors"
	"flag"
	// "fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "gw/_internal"

	"github.com/silenceper/log"
)

var (
	exit               chan bool // exit channel
	paths              []string  // watching paths
	cmdArgs            []string  // cmd args
	cmdName            string    // cmd name
	SubCommandErr      error     // error sub command
	CommandFormatError error     // error command format
	cfgFile            string    // config file
)

func init() {
	exit = make(chan bool)
	CommandFormatError = errors.New("Command err")
	SubCommandErr = errors.New("sub command invalid")
	initCommand()
}

func initCommand() {
	flag.StringVar(&cfgFile, "conf", "gowatch.yml", "default -conf=./gowatch.yml") // command
	flag.String("", "", "gowatch [ run/init ] [command] [args...]")                // command
}

func parseCommand() error {
	flag.Parse()
	args := flag.Args()
	// valid args
	if len(args) == 0 {
		return SubCommandErr
	}

	// get subcommand
	subCmd := args[0]

	switch subCmd {
	case "run":
		if len(args) < 2 {
			flag.Usage()
			return CommandFormatError
		}
		cmdName = args[1]
		if len(args) >= 3 {
			cmdArgs = args[2:]
		}
	case "init":
		// 输出到文件之后，结束程序
		log.Info("GoWactch Exit!")
		OutputDefaultConf("./gowatch.yml")
		os.Exit(2)
	default:
		flag.Usage()
		return SubCommandErr
	}

	return nil
}

/*
 * gowatch
 */
func main() {
	if parseCommand() != nil {
		return
	}
	w, err := NewWatcher()
	curPath, _ := os.Getwd()
	if err != nil {
		return
	}
	// parse config
	ParseConfig(cfgFile)
	cfg := GetInstance()
	log.Info(cfg.String())

	// set watcher
	paths = append(paths, curPath)
	paths = append(paths, cfg.ExternPaths...)
	WalkDirectoryRecursive(curPath, cfg.ExcludedPaths, &paths)
	AppendUnWatchRegexps(cfg.ExcludedRegexps...)
	StartWatch(w, paths, exit)

	// set dor
	InitDor(cmdName, cmdArgs, cfg.Envs)

	// handle os signal
	go HdlSignal()

	// main concurency keep going
	for {
		select {
		case <-exit:
			log.Info("GoWactch Exit!")
			Exit()
			os.Exit(2)
		}
		time.Sleep(3 * time.Second)
	}
}

// handle os signal
func HdlSignal() {
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
}
