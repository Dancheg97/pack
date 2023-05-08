package system

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"fmnx.io/core/pack/config"
)

func Callf(format string, a ...any) (string, error) {
	return Call(fmt.Sprintf(format, a...))
}

func Call(cmd string) (string, error) {
	if config.DebugMode {
		fmt.Println("Syscall - executing: ", cmd)
	}
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	if config.DebugMode {
		commad.Stdout = io.MultiWriter(&buf, os.Stdout)
		commad.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		commad.Stdout = &buf
		commad.Stderr = &buf
	}
	err := commad.Run()
	return buf.String(), err
}
