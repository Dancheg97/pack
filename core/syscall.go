package core

import (
	"fmt"
	"os"
	"os/exec"
)

func SystemCallf(format string, a ...any) error {
	return SystemCall(fmt.Sprintf(format, a...))
}

func SystemCall(cmd string) error {
	commad := exec.Command("bash", "-c", cmd)
	commad.Stdout = os.Stdout
	commad.Stderr = os.Stderr
	err := commad.Run()
	if err != nil {
		return err
	}
	return nil
}
