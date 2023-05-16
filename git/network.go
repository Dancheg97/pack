// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package git

// This package acts as library wrapper over git.
// You can use it to execute git calls.
// Package is safe for concurrent usage.

import (
	"errors"
	"fmt"
	"strings"

	"fmnx.su/core/pack/system"
)

// Clone git repository to specific folder.
func Clone(url string, dir string) error {
	out, err := system.Callf("git clone %s %s", url, dir)
	if err != nil {
		if !strings.Contains(out, "already exists and is not an empty dir") {
			return fmt.Errorf("git unable to clone: %s", out)
		}
	}
	return nil
}

// Get default branch for repo from remote.
func DefaultBranch(dir string) (string, error) {
	o, err := system.Callf("git -C %s remote show", dir)
	if err != nil {
		return ``, errors.New("git remote show error:\n" + o)
	}
	o = strings.Trim(o, "\n")
	o, err = system.Callf("git -C %s remote show %s", dir, o)
	if err != nil {
		return ``, errors.New("git remote show error:\n" + o)
	}
	rawInfo := strings.Split(o, "HEAD branch: ")[1]
	return strings.Split(rawInfo, "\n")[0], nil
}

// Get default branch for remote repository.
func LastCommitUrl(url string, branch string) (string, error) {
	o, err := system.Callf("git ls-remote -h %s", url)
	if err != nil {
		return ``, errors.New("git ls-remote error:\n" + o)
	}
	refs := strings.Split(strings.Trim(o, "\n"), "\n")
	for _, ref := range refs {
		if strings.HasSuffix(ref, branch) {
			return strings.Split(ref, "	")[0], nil
		}
	}
	return ``, errors.New("unable to find branch in remote repo")
}

// Pull changes for git repo under some branch in specified directory.
func Pull(dir string) error {
	_, err := system.Callf("git -C %s pull", dir)
	if err != nil {
		return errors.New("git pull error")
	}
	return nil
}
