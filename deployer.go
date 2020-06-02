package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	Commands []string
	Dir      string
	Log      *log.Logger
	ErrLog   *log.Logger
}

func NewRunner(commands []string) *Runner {
	return &Runner{
		Commands: commands,
		Log:      log.New(os.Stdout, "Log: ", log.LstdFlags),
		ErrLog:   log.New(os.Stdout, "Error: ", log.LstdFlags),
	}
}

func (r *Runner) run() error {
	for _, command := range r.Commands {
		if err := r.exec(command, r.Dir); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) exec(command string, dir string) error {
	c, args := r.parseCmd(command)
	cmd := exec.Command(c, args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		r.ErrLog.Println(c, err)
		return err
	}
	r.Log.Println(string(output))
	return nil
}

func (r *Runner) parseCmd(cmd string) (string, []string) {
	c := strings.Split(cmd, " ")
	return c[0], c[1:]
}
