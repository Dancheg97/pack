// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"fmt"
	"os"
	"strings"
)

// String template for adding new databases to pacman configuration
const dbtmpl = `
[%s]
Server = %s://%s
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
	HTTPS    bool
	TrustAll bool
}

// Add database to pacman configuration.
func AddConfigDatabase(p *RepositoryParameters) error {
	f, err := os.OpenFile(ConfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	link := p.Database + "/"
	p.Database = strings.Replace(p.Database, "/", ".", 1)
	trust := ``
	if p.TrustAll {
		trust = `SigLevel = Optional TrustAll`
	}
	protocol := `http`
	if p.HTTPS {
		protocol = `https`
	}
	tmpl := fmt.Sprintf(dbtmpl, p.Database, protocol, link, trust)
	_, err = f.Write([]byte(tmpl))
	return err
}
