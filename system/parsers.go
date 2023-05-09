// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io
package system

import (
	"fmt"
	"os"
	"strings"
)

// Swap some variable parameter in shell file and create substitute.
// Works only when it is single parameter declaration in shell file.
func SwapShellParameter(file string, param string, newval string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	splt1 := strings.Split(string(b), fmt.Sprintf("\n%s=", param))
	splt2 := strings.Split(splt1[1], "\n")
	join1 := strings.Join(splt2[1:], "\n")
	swp := fmt.Sprintf("\nswap%s=%s\n", param, splt2[0])
	rez := fmt.Sprintf("%s\n%s=%s%s%s", splt1[0], param, newval, swp, join1)
	rez = strings.ReplaceAll(rez, "$"+param, "$swap"+param)
	rez = strings.ReplaceAll(rez, "${"+param, "${swap"+param)
	rez = strings.ReplaceAll(rez, param+"(", "swap"+param+"(")
	return WriteFile(file, rez)
}

// Eject list parameters from shell file, typically PKGBUILD.
func GetShellList(file string, param string) ([]string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(string(f), param+"=(")
	if len(splitted) < 2 {
		return nil, nil
	}
	splitted = strings.Split(splitted[1], ")")
	dirtyParams := splitted[0]
	var cleanParams []string
	for _, param := range splitParams(dirtyParams) {
		cleanParams = append(cleanParams, cleanParameter(param))
	}
	return cleanParams, nil
}

func splitParams(params string) []string {
	// TODO rework add quotas check
	params = strings.ReplaceAll(params, "\n", " ")
	for strings.Contains(params, "  ") {
		params = strings.ReplaceAll(params, "  ", " ")
	}
	return strings.Split(strings.Trim(params, " "), " ")
}

func cleanParameter(param string) string {
	param = strings.ReplaceAll(param, "'", "")
	param = strings.ReplaceAll(param, "\"", "")
	return param
}
