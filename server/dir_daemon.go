// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"fmnx.su/core/pack/pacman"
	"github.com/radovskyb/watcher"
)

const (
	pkgext = ".pkg.tar.zst"
	dbext  = ".db.tar.gz"
)

// Parameters for directory db watcher.
type PkgDirParams struct {
	DbName    string
	WatchDir  string
	MkDirMode fs.FileMode

	InfoLogger Logger
	ErrLogger  Logger
}

var PkgDirDefault = PkgDirParams{
	DbName:     "localhost:4572",
	WatchDir:   "/var/cache/pacman/pkg",
	MkDirMode:  os.ModePerm,
	InfoLogger: log.Default(),
	ErrLogger:  log.Default(),
}

// This function is launching watcher for pacman cache directory, and constatly
// adding new arch packages to database in watched directory.
func PkgDirDaemon(p PkgDirParams) error {
	if err := os.MkdirAll(p.WatchDir, os.ModePerm); err != nil {
		return err
	}

	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move)
	if err := w.Add(p.WatchDir); err != nil {
		return err
	}

	go func() {
		for event := range w.Event {
			file := event.FileInfo.Name()
			if strings.HasSuffix(file, pkgext) {
				err := pacman.RepoAdd(
					path.Join(p.WatchDir, p.DbName+dbext),
					path.Join(p.WatchDir, file),
				)
				if err != nil {
					p.ErrLogger.Printf(
						"unable to add package %s to %s in %s",
						file, p.DbName, p.WatchDir,
					)
					continue
				}
				p.InfoLogger.Printf(
					"package %s added to db %s in dir %s",
					file, p.DbName, p.WatchDir,
				)
			}
		}
	}()
	return w.Start(time.Second)
}
