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

	. "github.com/yeqown/gw/_internal"

	"github.com/silenceper/log"
)

var (
	exit    = make(chan bool) // exit channel
	paths   []string          // watching paths
	cmdArgs []string          // cmd args
	cmdName string            // cmd name
	cfgFile string            // config file

	ErrSubCmd = errors.New("Command err")         // error sub command
	ErrCmdFmt = errors.New("sub command invalid") // error command format
)

func init() {
	flag.StringVar(&cfgFile, "gwconf", "gw.yml", "default -gwconf=./gw.yml") // command
	flag.String("", "", "gw [run/init] [command] [args...]")                 // command
}

// parse command line to do related command
func parseCommand() error {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return ErrSubCmd
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			flag.Usage()
			return ErrCmdFmt
		}
		cmdName = args[1]
		if len(args) >= 3 {
			cmdArgs = args[2:]
		}
	case "init":
		OutputDefaultConf("./gw.yml")
		log.Info("gw Exit!")
		os.Exit(2)
	default:
		flag.Usage()
		return ErrSubCmd
	}
	return nil
}

// the only entry
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

	log.Info("init dor")
	// set dor
	InitDor(cmdName, cmdArgs, cfg.Envs)

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
			log.Info("gw Exit!")
			Exit()
			os.Exit(2)
		}
		time.Sleep(3 * time.Second)
	}
}
