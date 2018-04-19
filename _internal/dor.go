/*
 * dor.go
 * process op, start, kill, restart, new
 */
package _internal

import (
	"github.com/silenceper/log"

	"fmt"
	"os"
	"os/exec"
	"syscall"
	// "time"
)

var (
	// EmptyArgs           = []string{} // default empty args
	// EmptyEnvs           = []string{} // default empty envs
	cmd          *exec.Cmd = nil // default cmd varible
	exitC        chan bool
	storeCmdArgs = []string{}
	storeCmdEnvs = []string{}
	storeCmdName = ""
)

func newCommand(cmdName string, cmdArgs, cmdEnvs []string) *exec.Cmd {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, cmdEnvs...)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd
}

func kill(cmd *exec.Cmd) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	if cmd == nil || cmd.Process == nil {
		return
	}

	if exited := <-exitC; exited {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGKILL)
		log.Infof("kill process success, %d", cmd.Process.Pid)
		return
	}
	panic(err)
}

func start(cmd *exec.Cmd) {
	cmd.Run()
	log.Infof("Command calling end: %s\n", cmd.ProcessState.String())
	exitC <- true
}

// final command will be like: "gowatch run ls -l"
// cmdArgs format: "", cmdEnv format: "GOOS=linux"
func InitDor(cmdName string, cmdArgs, cmdEnvs []string) {
	PATH := os.Getenv("path")
	cmdEnvs = append(cmdEnvs, fmt.Sprintf("%s=%s", "PATH", PATH))

	storeCmdArgs = cmdArgs
	storeCmdEnvs = cmdEnvs
	storeCmdName = cmdName

	cmd = newCommand(cmdName, cmdArgs, cmdEnvs)
	exitC = make(chan bool)
	log.Info("Command calling, please wait...")
	go start(cmd)
	exitC <- false
}

// hotReload one command
// if process has been killed, so renew one command
// else restart it
func hotReload() {
	kill(cmd)
	cmd = newCommand(storeCmdName, storeCmdArgs, storeCmdEnvs)
	go start(cmd)
	exitC <- false
}

// gowatch exit call this
func Exit() {
	if cmd.Process != nil && cmd.ProcessState != nil {
		kill(cmd)
	}
}
