package system

import (
	"os"
	"strings"
)

// Allows to eject parameters from PKGBUILD typycally to resolve dependencies.
func EjectShellList(file string, param string) ([]string, error) {
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
