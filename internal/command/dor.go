package command

/*
 * dor.go
 * process op, start, kill, restart, new
 */

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/silenceper/log"
)

// Command ...
type Command struct {
	exeCmd                     *exec.Cmd // default cmd varible
	exited                     bool
	storeCmdEnvs, storeCmdArgs []string
	storeCmdName               string
}

// New ...
func New(name string, args, envs []string) *Command {
	// PATH := os.Getenv("PATH")
	// log.Info(PATH)
	envs = append(envs, syscall.Environ()...)

	command := &Command{
		storeCmdArgs: args,
		storeCmdEnvs: envs,
		storeCmdName: name,
		exited:       false,
	}
	command.exeCmd = newExeCommand(name, args, envs)
	// log.Info("Command calling, please wait...")
	// go command.Start()
	// command.exit <- false
	return command
}

// Start ...
func (c *Command) Start() {
	go func() {
		c.exeCmd.Run()
		log.Infof("Command calling end: %s\n", c.exeCmd.ProcessState.String())
		c.exited = true
	}()
}

func (c *Command) quit() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	if c.exeCmd == nil || c.exeCmd.Process == nil {
		return
	}

	if c.exited {
		return
	}
	// resolved cannot kill all command by following:
	// syscall.Kill(-pgid, syscall.SIGKILL)
	pgid, err := syscall.Getpgid(c.exeCmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGKILL)
		log.Infof("kill process success, %d", c.exeCmd.Process.Pid)
		return
	}
	panic(err)
}

// HotReload ...
func (c *Command) HotReload() {
	c.quit()
	c.exeCmd = newExeCommand(c.storeCmdName, c.storeCmdArgs, c.storeCmdEnvs)
	c.Start()
	c.exited = false
}

// Exit ...
func (c *Command) Exit() {
	if c.exeCmd.Process != nil && c.exeCmd.ProcessState != nil {
		c.quit()
	}
}

func newExeCommand(name string, args, envs []string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, envs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}

// func kill(cmd *exec.Cmd) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Errorf("%s", err)
// 		}
// 	}()

// 	if cmd == nil || cmd.Process == nil {
// 		return
// 	}

// 	if exited := <-exitC; exited {
// 		return
// 	}
// 	// resolved cannot kill all command by following:
// 	// syscall.Kill(-pgid, syscall.SIGKILL)
// 	pgid, err := syscall.Getpgid(cmd.Process.Pid)
// 	if err == nil {
// 		syscall.Kill(-pgid, syscall.SIGKILL)
// 		log.Infof("kill process success, %d", cmd.Process.Pid)
// 		return
// 	}
// 	panic(err)
// }

// func start(cmd *exec.Cmd) {
// 	cmd.Run()
// 	log.Infof("Command calling end: %s\n", cmd.ProcessState.String())
// 	exitC <- true
// }

// // InitDor final command will be like: "gowatch run ls -l"
// // cmdArgs format: "", cmdEnv format: "GOOS=linux"
// func InitDor(cmdName string, cmdArgs, cmdEnvs []string) {
// 	// PATH := os.Getenv("PATH")
// 	// log.Info(PATH)
// 	cmdEnvs = append(cmdEnvs, syscall.Environ()...)

// 	storeCmdArgs = cmdArgs
// 	storeCmdEnvs = cmdEnvs
// 	storeCmdName = cmdName

// 	cmd = newExeCommand(cmdName, cmdArgs, cmdEnvs)
// 	exitC = make(chan bool)
// 	log.Info("Command calling, please wait...")
// 	go start(cmd)
// 	exitC <- false
// }

// // hotReload one command
// // if process has been killed, so renew one command
// // else restart it
// func hotReload() {
// 	kill(cmd)
// 	cmd = newExeCommand(storeCmdName, storeCmdArgs, storeCmdEnvs)
// 	go start(cmd)
// 	exitC <- false
// }

// // Exit gowatch exit call this
// func Exit() {
// 	if cmd.Process != nil && cmd.ProcessState != nil {
// 		kill(cmd)
// 	}
// }
