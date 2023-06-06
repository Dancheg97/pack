// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"errors"
	"io"
	"log"
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
type FsMirrorDaemon struct {
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

func (d *FsMirrorDaemon) init() error {
	if d.Link == "" {
		return errors.New("daemon mirror is not specified")
	}
	if d.Dir == "" {
		return errors.New("daemon load directory is not specified")
	}
	if d.ErrorLogger == nil {
		d.ErrorLogger = log.Default()
	}
	if d.InfoLogger == nil {
		d.InfoLogger = log.Default()
	}
	return nil
}

// Starts a daemon, loading files from remote file server and putting them into
// local directory. You should
func (d FsMirrorDaemon) Run() error {
	err := d.init()
	if err != nil {
		return err
	}

	for {
		client := http.Client{}

		resp, err := client.Get(d.Link)
		if err != nil {
			d.ErrorLogger.Printf("unable to get data from url: %s, %v", d.Link, err)
			return err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			d.ErrorLogger.Printf("unable to read resp: %s, %v", d.Link, err)
			return err
		}

		splt := strings.Split(string(b), href)
		if len(splt) < 2 {
			d.ErrorLogger.Printf("hrefs in mirror not found: %s", d.Link)
			return errors.New("unable to find packages with sign by link")
		}

		for _, packagefile := range splt {
			fileNoExtArr := strings.Split(packagefile, sigext+"\"")
			if len(fileNoExtArr) < 2 {
				continue
			}
			file := fileNoExtArr[0]

			if _, err := os.Stat(path.Join(d.Dir, file+pkgext)); err == nil {
				d.InfoLogger.Printf("package exists, skipping: %s", file)
				continue
			}

			err = d.LoadFile(file + pkgext)
			if err != nil {
				return err
			}

			err = d.LoadFile(file + sigext)
			if err != nil {
				return err
			}
			d.InfoLogger.Printf("package mirrored: %s", file)
		}

		time.Sleep(d.Dur)
	}
}

func (d FsMirrorDaemon) LoadFile(filename string) error {
	fileURL, err := url.Parse(d.Link + "/" + filename)
	if err != nil {
		d.ErrorLogger.Printf("unable to parse url: %s, %v", d.Link, err)
		return err
	}
	urlPath := fileURL.Path
	segments := strings.Split(urlPath, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	filepath := path.Join(d.Dir, fileName)
	file, err := os.Create(filepath)
	if err != nil {
		d.ErrorLogger.Printf("unable to create file: %s, %v", filepath, err)
		return err
	}
	client := http.Client{}

	// Put content on file
	resp, err := client.Get(d.Link + "/" + fileName)
	if err != nil {
		d.ErrorLogger.Printf("unable get file: %s, %v", d.Link, err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		d.ErrorLogger.Printf("unable to write to file: %s, %v", filepath, err)
		return err
	}
	defer file.Close()

	return nil
}
