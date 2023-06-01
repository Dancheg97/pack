// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// String template for adding new databases to pacman configuration
const dbtmpl = `
[%s]
Server = %s://%s/pacman
%s
`

var ConfigPath = "/etc/pacman.conf"

// Get list of connected databases in configuration.
func GetConfigDatabases() ([]string, error) {
	b, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}
	var dbs []string
	splt := strings.Split(string(b), "\n[")
	for _, line := range splt {
		dbs = append(dbs, strings.Split(line, "]")[0])
	}
	return dbs, nil
}

// Parameters for package repository, that will be used to connect to in
// pacman.conf file.
type RepositoryParameters struct {
	Database string
	HTTP     bool
	TrustAll bool
}

// Add database to pacman configuration.
func AddConfigDatabase(p *RepositoryParameters) error {
	fmt.Println(":: Adding database:", p.Database)
	protocol := "https"
	if p.HTTP {
		protocol = "http"
	}
	trust := ""
	if p.TrustAll {
		trust = "SigLevel = Optional TrustAll"
	}
	tmpl := fmt.Sprintf(dbtmpl, p.Database, protocol, p.Database, trust)
	for _, line := range strings.Split(tmpl, "\n") {
		cmd := exec.Command( //nolint:gosec
			"sudo", "bash", "-c",
			"echo "+line+" >> /etc/pacman.conf",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
