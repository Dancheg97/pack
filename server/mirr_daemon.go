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
	href   = "href=\""
)

// Parameters required to launch FS watcher for remote server mirroring.
type MirrFsParams struct {
	// Link to file server, from which files with .pkg.tar.zst and related .sig
	// files will be loaded. Use only trused mirrors!
	Link string
	// Directory, where resulting packages will be stored.
	Dir string
	// Syncronization duration.
	Dur time.Duration

	ErrorLogger Logger
	InfoLogger  Logger
}

// Starts a daemon, loading files from remote file server and putting them into
// local directory.
func MirrFsDaemon(p MirrFsParams) error {
	for {
		client := http.Client{}

		resp, err := client.Get(p.Link)
		if err != nil {
			p.ErrorLogger.Printf("unable to get data from url: %s, %v", p.Link, err)
			return err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			p.ErrorLogger.Printf("unable to read resp: %s, %v", p.Link, err)
			return err
		}

		splt := strings.Split(string(b), href)
		if len(splt) < 2 {
			p.ErrorLogger.Printf("hrefs in mirror not found: %s", p.Link)
			return errors.New("unable to find packages with sign by link")
		}

		for _, packagefile := range splt {
			fileNoExtArr := strings.Split(packagefile, sigext+"\"")
			if len(fileNoExtArr) < 2 {
				continue
			}
			file := fileNoExtArr[0]

			if _, err := os.Stat(path.Join(p.Dir, file+pkgext)); err == nil {
				p.InfoLogger.Printf("package exists, skipping: ", file)
				continue
			}

			err = LoadFile(LoadFileParams{
				Link:        p.Link + "/" + file + pkgext,
				Dir:         p.Dir,
				ErrorLogger: p.ErrorLogger,
				InfoLogger:  p.InfoLogger,
			})
			if err != nil {
				return err
			}

			err = LoadFile(LoadFileParams{
				Link:        p.Link + "/" + file + sigext,
				Dir:         p.Dir,
				ErrorLogger: p.ErrorLogger,
				InfoLogger:  p.InfoLogger,
			})
			if err != nil {
				return err
			}
			p.InfoLogger.Printf("package mirrored: ", file)
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

	ErrorLogger Logger
	InfoLogger  Logger
}

func LoadFile(p LoadFileParams) error {
	fileURL, err := url.Parse(p.Link)
	if err != nil {
		p.ErrorLogger.Printf("unable to parse url: %s, %v", p.Link, err)
		return err
	}
	urlPath := fileURL.Path
	segments := strings.Split(urlPath, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	filepath := path.Join(p.Dir, fileName)
	file, err := os.Create(filepath)
	if err != nil {
		p.ErrorLogger.Printf("unable to create file: %s, %v", filepath, err)
		return err
	}
	client := http.Client{}

	// Put content on file
	resp, err := client.Get(p.Link)
	if err != nil {
		p.ErrorLogger.Printf("unable get file: %s, %v", p.Link, err)
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		p.ErrorLogger.Printf("unable to write to file: %s, %v", filepath, err)
		return err
	}
	defer file.Close()

	p.InfoLogger.Printf("Downloaded a file %s with size %d", fileName, size)
	return nil
}
