// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pacman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	conf = "/etc/pacman.conf"
	tmpl = "cat <<EOF >> /etc/pacman.conf\n[%s]\nServer=https://%s/pack\nEOF"
)

// This function will read pacman.conf and add missing registries, that can
// returning the initial configuration string.
func AddRegistries(pkgs []string) (*string, error) {
	var registries []string
	for _, pkg := range pkgs {
		splt := strings.Split(pkg, "/")
		if len(splt) == 1 {
			continue
		}
		registries = append(registries, splt[0])
	}

	f, err := os.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	fstr := string(f)

	for _, registry := range registries {
		if !strings.Contains(fstr, fmt.Sprintf("\n[%s]\n", registry)) {
			addstr := fmt.Sprintf(tmpl, registry, registry)
			err := exec.Command("sudo", "bash", "-c", addstr).Run()
			if err != nil {
				return nil, errors.New("unable to add config with: " + addstr)
			}
		}
	}
	return &fstr, nil
}
