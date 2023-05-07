package system

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var Debug = false

func Callf(format string, a ...any) (string, error) {
	return Call(fmt.Sprintf(format, a...))
}

func Call(cmd string) (string, error) {
	if Debug {
		fmt.Println("Syscall - executing: ", cmd)
	}
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	if Debug {
		commad.Stdout = io.MultiWriter(&buf, os.Stdout)
		commad.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		commad.Stdout = &buf
		commad.Stderr = &buf
	}
	err := commad.Run()
	return buf.String(), err
}
