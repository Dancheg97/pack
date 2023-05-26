// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package pacman

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Query parameters for pacman packages.
type QueryOptions struct {
	// List packages explicitly installed. [--explicit]
	Explicit bool
	// List packages installed as dependencies. [--deps]
	Deps bool
	// Query only for packages installed from official repositories. [--native]
	Native bool
	// Query for packages installed from other sources. [--foreign]
	Foreign bool
	// Unrequired packages (not a dependency for other one). [--unrequired]
	Unrequired bool
	// Additional queue parameters.
	AdditionalParams []string
}

var QueryDefault = QueryOptions{}

type PackageInfo struct {
	Name    string
	Version string
}

// Get information about installed packages.
func Query(opts ...QueryOptions) ([]PackageInfo, error) {
	o := formOptions(opts, &QueryDefault)

	args := []string{"-Q"}
	if o.Explicit {
		args = append(args, "--explicit")
	}
	if o.Deps {
		args = append(args, "--deps")
	}
	if o.Native {
		args = append(args, "--native")
	}
	if o.Foreign {
		args = append(args, "--foreign")
	}
	if o.Unrequired {
		args = append(args, "--unrequired")
	}
	args = append(args, o.AdditionalParams...)

	var b bytes.Buffer
	cmd := exec.Command(pacman, args...)
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		if b.String() == "" {
			return nil, nil
		}
		return nil, errors.New("unable to query packages: " + b.String())
	}
	return parsePackages(b.String()), nil
}

func parsePackages(q string) []PackageInfo {
	var rez []PackageInfo
	for _, v := range strings.Split(q, "\n") {
		if v == "" {
			break
		}
		splt := strings.Split(v, " ")
		rez = append(rez, PackageInfo{
			Name:    splt[0],
			Version: splt[1],
		})
	}
	return rez
}

type PackageInfoFull struct {
	Name          string
	Version       string
	Description   string
	Architecture  string
	Url           string
	Licenses      string
	Groups        string
	Provides      string
	DependsOn     string
	OptionalDeps  string
	RequiredBy    string
	OptionalFor   string
	ConflictsWith string
	Replaces      string
	InstalledSize string
	Packager      string
	BuildDate     string
	InstallDate   string
	InstallReason string
	InstallScript string
	ValidatedBy   string
}

// Get info about package.
func Info(pkg string) (*PackageInfoFull, error) {
	var b bytes.Buffer
	cmd := exec.Command(pacman, "-Qi", pkg)
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		return nil, errors.New("unable to get info: " + b.String())
	}
	out := b.String()

	return &PackageInfoFull{
		Name:          parseField(out, "Name            : "),
		Version:       parseField(out, "Version         : "),
		Description:   parseField(out, "Description     : "),
		Architecture:  parseField(out, "Architecture    : "),
		Url:           parseField(out, "URL             : "),
		Licenses:      parseField(out, "Licenses        : "),
		Groups:        parseField(out, "Groups          : "),
		Provides:      parseField(out, "Provides        : "),
		DependsOn:     parseField(out, "Depends On      : "),
		OptionalDeps:  parseField(out, "Optional Deps   : "),
		RequiredBy:    parseField(out, "Required By     : "),
		OptionalFor:   parseField(out, "Optional For    : "),
		ConflictsWith: parseField(out, "Conflicts With  : "),
		Replaces:      parseField(out, "Replaces        : "),
		InstalledSize: parseField(out, "Installed Size  : "),
		Packager:      parseField(out, "Packager        : "),
		BuildDate:     parseField(out, "Build Date      : "),
		InstallDate:   parseField(out, "Install Date    : "),
		InstallReason: parseField(out, "Install Reason  : "),
		InstallScript: parseField(out, "Install Script  : "),
		ValidatedBy:   parseField(out, "Validated By    : "),
	}, nil
}

func parseField(full string, field string) string {
	splt := strings.Split(full, field)
	return strings.Split(splt[1], "\n")[0]
}

// Outdated package.
type OutdatedPackage struct {
	Name           string
	CurrentVersion string
	NewVersion     string
}

// Get information about outdated packages.
func Outdated() ([]OutdatedPackage, error) {
	var b bytes.Buffer
	cmd := exec.Command(pacman, "-Qu")
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		if b.String() == `` {
			return nil, nil
		}
		return nil, errors.New("unable to get info: " + b.String())
	}
	out := b.String()
	return parseOutdated(out), nil
}

func parseOutdated(o string) []OutdatedPackage {
	var rez []OutdatedPackage
	for _, line := range strings.Split(o, "\n") {
		if line == "" {
			break
		}
		splt := strings.Split(line, " ")
		rez = append(rez, OutdatedPackage{
			Name:           splt[0],
			CurrentVersion: splt[1],
			NewVersion:     splt[3],
		})
	}
	return rez
}
