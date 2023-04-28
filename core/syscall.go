package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func SystemCallf(format string, a ...any) error {
	return SystemCall(fmt.Sprintf(format, a...))
}

func SystemCall(cmd string) error {
	fmt.Println("Running: ", cmd)
	commad := exec.Command("bash", "-c", cmd)
	commad.Stdout = os.Stdout
	commad.Stderr = os.Stderr
	err := commad.Run()
	if err != nil {
		return err
	}
	return nil
}
func SystemCallOutf(format string, a ...any) (string, error) {
	return SystemCallOut(fmt.Sprintf(format, a...))
}

func SystemCallOut(cmd string) (string, error) {
	fmt.Println("Running: ", cmd)
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	commad.Stdout = &buf
	commad.Stderr = &buf
	err := commad.Run()
	return buf.String(), err
}
