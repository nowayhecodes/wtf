package shell

import (
	"os"
	"os/exec"
	"strings"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(cmd string) error {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return nil
	}

	command := exec.Command(args[0], args[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
} 