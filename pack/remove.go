// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"fmnx.su/core/pack/msgs"
	"fmnx.su/core/pack/pacman"
)

type RemoveParameters struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader

	// Ask for confirmation when deleting package.
	Confirm bool
	// Leave package dependencies in the system (removed by default).
	Norecursive bool
	// Leave package configs in the system (removed by default).
	Nocfgs bool
	// Remove packages and all packages that depend on them.
	Cascade bool
	// Custom distribution name that will be used for package deletion.
	Distro string
	// Use insecure connection for remote deletions.
	Insecure bool
}

func removeDefault() *RemoveParameters {
	return &RemoveParameters{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Distro: "archlinux",
	}
}

func Remove(args []string, prms ...RemoveParameters) error {
	p := formOptions(prms, removeDefault)

	local, remote := splitRemoved(args)

	if len(local) > 0 {
		err := pacman.RemoveList(local, pacman.RemoveParameters{
			Sudo:        true,
			NoConfirm:   !p.Confirm,
			Recursive:   !p.Norecursive,
			WithConfigs: !p.Nocfgs,
			Cascade:     p.Cascade,
			Stdout:      p.Stdout,
			Stderr:      p.Stderr,
			Stdin:       p.Stdin,
		})
		if err != nil {
			return err
		}
	}

	if len(remote) > 0 {
		email, err := gnupgEmail()
		if err != nil {
			return err
		}
		if p.Distro == "" {
			p.Distro = "archlinux"
		}
		msgs.Amsg(p.Stdout, "Removing remote packages as "+email)
		for i, pkg := range remote {
			msgs.Smsg(p.Stdout, "Removing "+pkg, i, len(remote))
			err := rmRemote(p, pkg, email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Splits packages that will be removed locally and on remote.
func splitRemoved(pkgs []string) ([]string, []string) {
	var local []string
	var remote []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "/") {
			remote = append(remote, pkg)
			continue
		}
		local = append(local, pkg)
	}
	return local, remote
}

func splitPkg(pkg string) (string, string, string) {
	splt := strings.Split(pkg, "/")
	if len(splt) == 2 {
		return splt[0], ``, splt[1]
	}
	return splt[0], splt[1], splt[2]
}

// Function that will be used to remove remote package.
func rmRemote(p *RemoveParameters, pkg, email string) error {
	t := time.Now().Format(time.RFC3339)

	remote, owner, target := splitPkg(pkg)

	err := os.WriteFile("packdel", []byte(t+remote+owner+target), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll("packdel")

	cmd := exec.Command("gpg", "--sign", "packdel")
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	sigdata, err := os.ReadFile("packdel")
	if err != nil {
		return err
	}

	prfx := "https://"
	if p.Insecure {
		prfx = "http://"
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		prfx+path.Join(remote, "api/packages", owner, "arch/remove"),
		bytes.NewReader(sigdata),
	)
	if err != nil {
		return err
	}

	req.Header.Add("email", email)
	req.Header.Add("distro", p.Distro)
	req.Header.Add("target", target)

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(err, errors.New(resp.Status))
		}
		return fmt.Errorf("%s, %s %s", resp.Status, string(b), pkg)
	}
	return nil
}
