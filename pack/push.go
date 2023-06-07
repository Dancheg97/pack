// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

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

	"fmnx.su/core/pack/tmpl"
)

// Parameters that will be used to execute push command.
type PushParameters struct {
	// Directory to read package files and signatures.
	Directory string
	// Which protocol to use for connection.
	Protocol string
	// Custom endpoint for package push
	Endpoint string
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
}

func pushdefault() *PushParameters {
	return &PushParameters{
		Protocol:  "https",
		Endpoint:  "/api/pack/push",
		Directory: "/var/cache/pacman/pkg",
	}
}

// Push your package to registry.
func Push(args []string, prms ...PushParameters) error {
	p := formOptions(prms, pushdefault)

	email, err := gnupgEmail(p.Stderr)
	if err != nil {
		return err
	}

	pkgs, _, err := formatpkgs(args, p.Stderr)
	if err != nil {
		return err
	}

	err = checkRegistries(pkgs, p.Stderr)
	if err != nil {
		return err
	}

	filenames, err := listPkgFilenames(p.Directory, p.Stderr)
	if err != nil {
		return err
	}

	pprms, err := fillfileinfo(fillparams{
		filenames: filenames,
		packages:  pkgs,
		directory: p.Directory,
		ew:        p.Stderr,
	})
	if err != nil {
		return err
	}

	for _, pp := range pprms {
		err = push(pp, email, p.Protocol, p.Endpoint)
		if err != nil {
			return err
		}
	}

	return nil
}

// This function will be used to get email from user's GnuPG identitry.
func gnupgEmail(ew io.Writer) (string, error) {
	gnupgident, err := gnuPGIdentity()
	if err != nil {
		ew.Write([]byte(tmpl.ErrEmailRead + err.Error())) //nolint
		return ``, err
	}
	return strings.ReplaceAll(strings.Split(gnupgident, "<")[1], ">", ""), nil
}

// Check if all packages have registries where they will be pushed to.
func checkRegistries(pkgs []RegistryPkg, ew io.Writer) error {
	for _, pkg := range pkgs {
		if pkg.Registry == "" {
			ew.Write([]byte(tmpl.ErrNoRegistry + pkg.Name)) //nolint
			return errors.New("no registry for package " + pkg.Name)
		}
	}
	return nil
}

// Structure including base registry parameters and information about file
// pathes requied to push packages.
type PushPkg struct {
	RegistryPkg
	// Name of the file which will be pushed.
	Filename string
	// Path to file which will be read and pushed.
	PkgPath string
	// Signature encoded to base64 string to check.
	Signature string
}

// List file names in provided cache directory.
func listPkgFilenames(dir string, ew io.Writer) ([]string, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		ew.Write([]byte(tmpl.ErrDirRead + dir + " " + err.Error())) //nolint
		return nil, err
	}
	var fns []string
	for _, de := range des {
		fn := de.Name()
		if strings.HasSuffix(fn, ".pkg.tar.zst") {
			fns = append(fns, fn)
		}
	}
	return fns, nil
}

type fillparams struct {
	filenames []string
	packages  []RegistryPkg
	directory string
	ew        io.Writer
}

// Create array of package arguements, that will be pushed to registry.
func fillfileinfo(p fillparams) ([]PushPkg, error) {
	var ppkgs []PushPkg
	for _, pkg := range p.packages {
		for i, filename := range p.filenames {
			if !strings.Contains(filename, pkg.Name) {
				continue
			}
			pkgname, err := ejectpkgname(filename, p.ew)
			if err != nil {
				return nil, err
			}
			if pkgname == pkg.Name {
				pkgpath := path.Join(p.directory, filename)
				sigbase64, err := readpkgsign(pkgpath+".sig", p.ew)
				if err != nil {
					return nil, err
				}
				ppkgs = append(ppkgs, PushPkg{
					RegistryPkg: pkg,
					Filename:    filename,
					PkgPath:     pkgpath,
					Signature:   sigbase64,
				})
				break
			}
			if i == len(p.filenames) {
				p.ew.Write([]byte(tmpl.ErrPkgNotFound + pkg.Name)) //nolint
				return nil, errors.New("package file not found " + pkg.Name)
			}
		}
	}
	return ppkgs, nil
}

// Eject package name from file name.
func ejectpkgname(filename string, ew io.Writer) (string, error) {
	pkgsplt := strings.Split(filename, "-")
	if len(pkgsplt) < 4 {
		ew.Write([]byte(tmpl.ErrBrokenPkgFile + filename)) //nolint
		return ``, errors.New("invalid package in cache: " + filename)
	}
	return strings.Join(pkgsplt[:len(pkgsplt)-3], "-"), nil
}

// Read package signature and encode to base64.
func readpkgsign(path string, ew io.Writer) (string, error) {
	err := exec.Command("bash", "-c", "sudo chmod 0777 "+path).Run() //nolint
	if err != nil {
		ew.Write([]byte(tmpl.ErrSigRead + path)) //nolint
		return ``, err
	}
	sigbytes, err := os.ReadFile(path)
	if err != nil {
		ew.Write([]byte(tmpl.ErrSigRead + path)) //nolint
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(sigbytes), nil
}

// This function pushes package to registry via http.
func push(p PushPkg, email string, protocol string, endpoint string) error {
	packagefile, err := os.Open(p.PkgPath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		protocol+"://"+p.Registry+endpoint,
		packagefile,
	)
	if err != nil {
		return err
	}

	req.Header.Add("file", p.Filename)
	req.Header.Add("email", email)
	req.Header.Add("sign", p.Signature)
	req.Header.Add("owner", p.Owner)

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

	fmt.Println("[PUSH] - package delivered: " + p.Name)
	return nil
}
