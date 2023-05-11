// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package system

import (
	"os"
	"strings"
)

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
