// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"fmnx.su/core/pack/pacman"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{"p"},
	Short:   "📨 push packages",
	Long: `📨 push packages

This command will find package and sign in package cache directory and push it
to provided registry. Registry should be included with package name, example:

pack p localhost:4572/linux-zen`,
	Run: Push,
}

func Push(cmd *cobra.Command, args []string) {
	gnupgident, err := pacman.GetGnupgIdentity()
	CheckErr(err)
	var pkgs []*Package
	for _, pkg := range args {
		p, err := FormPackage(pkg)
		p.Email = strings.ReplaceAll(strings.Split(gnupgident, "<")[1], ">", "")
		CheckErr(err)
		pkgs = append(pkgs, p)
	}
	for _, pkg := range pkgs {
		err := PushPkg(pkg)
		CheckErr(err)
	}
}

type Package struct {
	Registry string
	PkgName  string
	Filename string
	PkgFile  string
	SigFile  string
	Email    string
}

// This function will find the latest version of package in cache direcotry and
// then push it to registry specified in package name provided in argiement.
func FormPackage(pkg string) (*Package, error) {
	splt := strings.Split(pkg, "/")
	if len(splt) != 2 {
		msg := "error: package should contain registry and name: "
		return nil, errors.New(msg + pkg)
	}
	registry := splt[0]
	pkgname := splt[1]
	des, err := os.ReadDir(pacmancache)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		filename := de.Name()
		if !strings.HasSuffix(filename, pkgext) {
			continue
		}
		pkgsplt := strings.Split(filename, "-")
		if len(pkgsplt) < 4 {
			return nil, errors.New("invalid package in cache: " + filename)
		}
		namecheck := strings.Join(pkgsplt[:len(pkgsplt)-3], "-")
		if pkgname == namecheck {
			return &Package{
				Registry: registry,
				PkgName:  pkgname,
				Filename: filename,
				PkgFile:  path.Join(pacmancache, filename),
				SigFile:  path.Join(pacmancache, filename+".sig"),
			}, err
		}
	}
	return nil, errors.New("package not found in cache: " + pkgname)
}

// This function pushes package to registry via http.
func PushPkg(p *Package) error {
	packagefile, err := os.Open(p.PkgFile)
	if err != nil {
		return err
	}
	fmt.Println(":: Retrieving package signature access.")
	err = exec.Command( // nolint:gosec
		"bash", "-c",
		"sudo chmod 0777 "+p.SigFile,
	).Run()
	if err != nil {
		return err
	}
	sigbytes, err := os.ReadFile(p.SigFile)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		protocol+p.Registry+pushendpoint,
		packagefile,
	)
	if err != nil {
		return err
	}

	req.Header.Add(file, p.Filename)
	req.Header.Add(email, p.Email)
	req.Header.Add(sign, base64.StdEncoding.EncodeToString(sigbytes))

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
		return errors.New(resp.Status + " " + string(b))
	}

	fmt.Println("[PUSH] - package delivered: " + p.PkgFile)
	return nil
}
