/*
 * 执行者
 */
package internal

import (
	"github.com/silenceper/log"

	// "io"
	"os"
	"os/exec"
)

var (
	cmd       *exec.Cmd = nil        // default cmd varible
	EmptyArgs           = []string{} // default empty args
	EmptyEnvs           = []string{} // default empty envs

	storeCmdArgs = []string{}
	storeCmdEnvs = []string{}
	storeCmdName = ""
)

func newCommand(cmdName string, cmdArgs, cmdEnvs []string) *exec.Cmd {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, cmdEnvs...)
	return cmd
}

func kill(cmd *exec.Cmd) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			panic(err)
		}
	}
	return
}

func start(cmd *exec.Cmd) {
	log.Info("Command Calling")
	go cmd.Run()
}

// cmdArgs format: "", cmdEnv format: "GOOS=linux"
func InitDor(cmdName string, cmdArgs, cmdEnvs []string) {
	// store
	storeCmdArgs = cmdArgs
	storeCmdEnvs = cmdEnvs
	storeCmdName = cmdName
	// may need valid args and env input
	cmd = newCommand(cmdName, cmdArgs, cmdEnvs)
	// init start cmd
	start(cmd)
}

func hotReload() {
	if !cmd.ProcessState.Exited() {
		kill(cmd)
	} else {
		cmd = newCommand(storeCmdName, storeCmdArgs, storeCmdEnvs)
	}
	start(cmd)
}
