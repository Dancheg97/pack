// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

// Optional parameters for pacman sync command.
type SyncOptions struct {
	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Do not reinstall up to date packages. [--needed]
	Needed bool
	// Do not ask for any confirmation. [--noconfirm]
	NoConfirm bool
	// Do not show a progress bar when downloading files. [--noprogressbar]
	NoProgressBar bool
	// Do not execute the install scriptlet if one exists. [--noscriptlet]
	NoScriptlet bool
	// Install packages as non-explicitly installed. [--asdeps]
	AsDeps bool
	// Install packages as explictly installed. [--asexplict]
	AsExplict bool
	// Download fresh package databases from the server. [--refresh]
	Refresh bool
	// Upgrade programms that are outdated. [--sysupgrade]
	Upgrade bool
	// Only download, but do not install package. [--downloadonly]
	DownloadOnly bool
	// Clean old packages from cache directory. [--clean]
	Clean bool
	// Clean all packages from cache directory. [-cc]
	CleanAll bool
	// Where command will write output text.
	Stdout io.Writer
	// Where command will write output text.
	Stderr io.Writer
	// Stdin from user is command will ask for something.
	Stdin io.Reader
	// Additional parameters, that will be appended to command as arguements.
	AdditionalParams []string
}

// Those are options that will be set up by default on program execution.
var SyncDefault = SyncOptions{
	Sudo:      true,
	Needed:    true,
	NoConfirm: true,
	Stdout:    os.Stdout,
	Stderr:    os.Stderr,
}

// Executes pacman sync command. This command will read sync options and form
// command based on first elements from the array.
func Sync(pkgs string, opts ...SyncOptions) error {
	return SyncList(strings.Split(pkgs, " "), opts...)
}

// Sync command for package string list.
func SyncList(pkgs []string, opts ...SyncOptions) error {
	o := formOptions(opts, &SyncDefault)

	args := []string{"-S"}
	if o.Needed {
		args = append(args, "--needed")
	}
	if o.NoConfirm {
		args = append(args, "--noconfirm")
	}
	if o.NoProgressBar {
		args = append(args, "--noprogressbar")
	}
	if o.NoScriptlet {
		args = append(args, "--noscriptlet")
	}
	if o.AsDeps {
		args = append(args, "--asdeps")
	}
	if o.AsExplict {
		args = append(args, "--asexplicit")
	}
	if o.Refresh {
		args = append(args, "--refresh")
	}
	if o.Upgrade {
		args = append(args, "--sysupgrade")
	}
	if o.DownloadOnly {
		args = append(args, "--downloadonly")
	}
	if o.DownloadOnly {
		args = append(args, "--downloadonly")
	}
	if o.Clean {
		args = append(args, "--clean")
	}
	if o.CleanAll {
		args = append(args, "-cc")
	}
	args = append(args, o.AdditionalParams...)
	args = append(args, pkgs...)

	cmd := SudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Stdin = o.Stdin

	mu.Lock()
	defer mu.Unlock()
	return cmd.Run()
}

// Options to apply when searching for some package.
type SearchOptions struct {
	// Run with sudo priveleges. [sudo]
	Sudo bool
	// Download fresh package databases from the server. [--refresh]
	Refresh bool
	// Stdin from user is command will ask for something.
	Stdin io.Reader
}

// Structure to recieve from search result
type SearchResult struct {
	Repo    string
	Name    string
	Version string
	Desc    string
}

var SearchDefault = SearchOptions{
	Sudo:    true,
	Refresh: true,
	Stdin:   os.Stdin,
}

// Search for packages.
func Search(re string, opts ...SearchOptions) ([]SearchResult, error) {
	o := formOptions(opts, &SearchDefault)

	args := []string{"-Ss"}
	if o.Refresh {
		args = append(args, "--refresh")
	}
	args = append(args, re)

	var b bytes.Buffer
	cmd := SudoCommand(o.Sudo, pacman, args...)
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.Stdin = os.Stdin

	mu.Lock()
	defer mu.Unlock()
	err := cmd.Run()

	if err != nil {
		if b.String() == `` {
			return nil, nil
		}
		return nil, errors.New("unable to search: " + b.String())
	}
	return serializeOutput(b.String()), nil
}

func serializeOutput(output string) []SearchResult {
	if strings.HasPrefix(output, ":: Synchronizing package databases") {
		splt := strings.Split(output, "downloading...\n")
		output = splt[len(splt)-1]
	}
	var rez []SearchResult
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if line == `` {
			break
		}
		if i%2 == 1 {
			continue
		}
		splt := strings.Split(line, " ")
		repoName := strings.Split(splt[0], "/")
		rez = append(rez, SearchResult{
			Repo:    repoName[0],
			Name:    repoName[1],
			Version: splt[1],
			Desc:    strings.Trim(lines[i+1], " "),
		})
	}
	return rez
}
