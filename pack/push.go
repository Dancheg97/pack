// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// This command will push provided package to pack server using existing pack
// credentials, otherwise it will ask to create them.
// func Push(pkgs ...string) error {
// 	u, p, err := getCreads()
// 	if err != nil {
// 		return err
// 	}
// 	for _, pkg := range pkgs {
// 		fmt.Println(":: Pushing package", pkg)
// 		os.read
// 	}
// 	return nil
// }

func getCreads() (string, string, error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	b, err := os.ReadFile(path.Join(u.HomeDir, ".packcfg"))
	if err != nil {
		fmt.Println(":: Unable to find pack credentials.")
		n, p, err := askForCreds()
		if err != nil {
			return "", "", err
		}
		err = os.WriteFile(".packcfg", []byte(n+" "+p), 0644)
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
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
