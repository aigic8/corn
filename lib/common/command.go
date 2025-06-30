package common

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// TODO: handle stdout and stderr better
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	log.Printf("command '%s' started", JoinCommandParts(command, args))

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
