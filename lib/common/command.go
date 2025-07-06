package common

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type RunCommandOpts struct {
	Cmd    string
	Args   []string
	Stdout io.Writer
	Stderr io.Writer
}

func RunCommand(o *RunCommandOpts) error {
	cmd := exec.Command(o.Cmd, o.Args...)

	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command exited with error: %w", err)
	}
	return nil
}

func SeparateCommandFromArgs(raw string) (string, []string) {
	sections := strings.Split(raw, " ")
	command := sections[0]

	if len(sections) < 2 {
		return command, []string{}
	}
	return command, sections[1:]
}

func JoinCommandParts(command string, args []string) string {
	return command + " " + strings.Join(args, " ")
}
