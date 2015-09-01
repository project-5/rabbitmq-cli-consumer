package command

import (
	"log"
	"io"
	"os"
	"os/exec"
)

type CommandExecuter struct {
	errLogger *log.Logger
	infLogger *log.Logger
}

func New(errLogger, infLogger *log.Logger) *CommandExecuter {
	return &CommandExecuter{
		errLogger: errLogger,
		infLogger: infLogger,
	}
}

func (me CommandExecuter) Execute(cmd *exec.Cmd) bool {
	me.infLogger.Println("Processing message...")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stdout, os.Stderr)
	
	err := cmd.Run()

	if err != nil {
		me.errLogger.Printf("Processing error: %s\n", err)
		return false
	}

	me.infLogger.Println("Processed!")

	return true
}
