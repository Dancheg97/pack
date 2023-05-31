// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// This command will push provided package to pack server using existing pack
// credentials, otherwise it will ask to create them.
func Push(repo string, pkgs ...string) error {
	u, p, err := getCreads()
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		fmt.Println(":: Pushing package", pkg)
		des, err := os.ReadDir("/var/cache/pacman/pkg")
		if err != nil {
			return err
		}
		for _, de := range des {
			if strings.HasSuffix(de.Name(), ".pkg.tar.zst") {
				splt := strings.Split(de.Name(), "-")
				filepkg := strings.Join(splt[0:len(splt)-3], "-")
				if pkg == filepkg {
					f, err := os.Open(de.Name())
					if err != nil {
						return err
					}
					r, err := http.NewRequest(
						"POST", repo+"/pacman/push", f,
					)
					if err != nil {
						return err
					}
					r.Header.Set("user", u)
					r.Header.Set("password", p)
					r.Header.Set("file", de.Name())
					c := http.Client{}
					_, err = c.Do(r)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func getCreads() (string, string, error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	cfgpath := path.Join(u.HomeDir, ".packcfg")
	b, err := os.ReadFile(cfgpath)
	if err != nil {
		fmt.Println(":: Unable to find pack credentials.")
		n, p, err := askForCreds()
		if err != nil {
			return "", "", err
		}
		err = os.WriteFile(cfgpath, []byte(n+" "+p), 0600)
		if err != nil {
			return "", "", err
		}
		return n, p, nil
	}
	splt := strings.Split(string(b), " ")
	return splt[0], splt[1], nil
}

func askForCreds() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
