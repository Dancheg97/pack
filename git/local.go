// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package git

// This package acts as library wrapper over git.
// You can use it to execute git calls.

import (
	"errors"
	"strings"

	"fmnx.io/core/pack/system"
)

// Switch repo to branch/tag/commit.
func Checkout(dir string, target string) error {
	o, err := system.Callf("git -C %s checkout %s ", dir, target)
	if err != nil {
		if !strings.HasPrefix(o, "Already on ") {
			return errors.New("git unable to find checkout target:\n" + o)
		}
	}
	return nil
}

// Clean git repository - all changes in tracked files, newly created files and
// files under gitignore.
func Clean(dir string) error {
	o, err := system.Callf("git -C %s clean -xdf", dir)
	if err != nil {
		return errors.New("git unable to clean xdf:\n" + o)
	}
	o, err = system.Callf("git -C %s reset --hard", dir)
	if err != nil {
		return errors.New("git unable to reset -hard:\n" + o)
	}
	return nil
}

// Get last commit hash for git repo in a branch.
func LastCommitDir(dir string, branch string) (string, error) {
	command := `git -C ` + dir + ` log -n 1 --pretty=format:"%H" ` + branch
	o, err := system.Call(command)
	if err != nil {
		return ``, errors.New("git unable to log:\n" + o)
	}
	return strings.Trim(o, "\n"), nil
}

// Get git installation url and convert it to https format.
func Url(dir string) (string, error) {
	out, err := system.Callf("git -C %s config --get remote.origin.url", dir)
	if err != nil {
		return ``, errors.New("git unable to get remote url:\n" + out)
	}
	out = strings.Trim(out, "\n")
	out = strings.Replace(out, "git@", "https://", 1)
	out = strings.Replace(out, ":", "/", 1)
	out = strings.Replace(out, ".git", "", 1)
	return out, nil
}
