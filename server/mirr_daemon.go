// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"os"
	"strings"
)

const (
	sigext = ".pkg.tar.zst.sig"
)

// Parameters required to launch FS watcher for remote server mirroring.
type MirrFsParams struct {
	// Link to file server, from which files with .pkg.tar.zst and related .sig
	// files will be loaded. Use only trused mirrors!
	Link string
	// Directory, where resulting packages will be stored.
	Dir string
	// Logger, that will be used to log errors.
	Logger Logger
	// Syncronization duration.
	Dur time.Duration
}

// Starts a daemon, loading files from remote file server and putting them into
// local directory.
func MirrFsDaemon(p MirrFsParams) error {
	for {
		client := http.Client{}

		resp, err := client.Get(p.Link)
		if err != nil {
			p.Logger.Printf("unable to get data from url: %s, %v", p.Link, err)
			return err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			p.Logger.Printf("unable to read resp: %s, %v", p.Link, err)
			return err
		}

		splt := strings.Split(string(b), sigext)
		if len(splt) < 2 {
			p.Logger.Printf("signed packages not found: %s", p.Link)
			return errors.New("unable to find packages with sign by link")
		}

		for i, v := range splt {
			if i+1 == len(splt) {
				break
			}
			quotesplt := strings.Split(v, "\"")
			nlsplt := strings.Split(quotesplt[len(quotesplt)-1], "\n")

			pkglink := p.Link + nlsplt[len(nlsplt)-1] + pkgext
			err = LoadFile(LoadFileParams{
				Link:   pkglink,
				Dir:    p.Dir,
				Logger: p.Logger,
			})
			if err != nil {
				return err
			}

			siglink := p.Link + nlsplt[len(nlsplt)-1] + sigext
			err = LoadFile(LoadFileParams{
				Link:   siglink,
				Dir:    p.Dir,
				Logger: p.Logger,
			})
			if err != nil {
				return err
			}
		}

		time.Sleep(p.Dur)
	}
}

// Parameters for single file download from remote.
type LoadFileParams struct {
	// Link to downloaded file.
	Link string
	// Directory, where resulting packages will be stored.
	Dir string
	// Logger, that will be used to log errors.
	Logger Logger
}

func LoadFile(p LoadFileParams) error {
	fileURL, err := url.Parse(p.Link)
	if err != nil {
		p.Logger.Printf("unable to parse url: %s, %v", p.Link, err)
		return err
	}
	urlPath := fileURL.Path
	segments := strings.Split(urlPath, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	filepath := path.Join(p.Dir, fileName)
	file, err := os.Create(filepath)
	if err != nil {
		p.Logger.Printf("unable to create file: %s, %v", filepath, err)
		return err
	}
	client := http.Client{}

	// Put content on file
	resp, err := client.Get(p.Link)
	if err != nil {
		p.Logger.Printf("unable get file: %s, %v", p.Link, err)
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		p.Logger.Printf("unable to write to file: %s, %v", filepath, err)
		return err
	}
	defer file.Close()

	p.Logger.Printf("Downloaded a file %s with size %d", fileName, size)
	return nil
}
