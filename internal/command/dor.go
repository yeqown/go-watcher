package command

/*
 * dor.go
 * process op, start, kill, restart, new
 */

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/yeqown/go-watcher/internal/log"
)

// Command ... contains command context of 'name', 'args', and 'envs'
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
	return command
}

// Start ... execute the command
func (c *Command) Start() {
	go func() {
		c.exeCmd.Run()
		// state is: %s\n", c.exeCmd.ProcessState.String()
		log.Info("command executed done!")
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
	// [resolved] TOFIX: cannot kill all command by following:
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
