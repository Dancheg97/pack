package system

import (
	"os"
	"strings"
)

func EjectShList(file string, param string) ([]string, error) {
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
	for _, param := range SplitParams(dirtyParams) {
		cleanParams = append(cleanParams, CleanParameter(param))
	}
	return cleanParams, nil
}

func SplitParams(params string) []string {
	params = strings.ReplaceAll(params, "\n", " ")
	for strings.Contains(params, "  ") {
		params = strings.ReplaceAll(params, "  ", " ")
	}
	return strings.Split(strings.Trim(params, " "), " ")
}

func CleanParameter(param string) string {
	param = strings.ReplaceAll(param, "'", "")
	param = strings.ReplaceAll(param, "\"", "")
	param = strings.Split(param, "=")[0]
	param = strings.Split(param, ">")[0]
	return param
}
