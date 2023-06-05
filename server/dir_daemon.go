// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

import (
	"bytes"
	"errors"
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
type PkgDirDaemon struct {
	DbName    string
	WatchDir  string
	MkDirMode fs.FileMode

	InfoLogger Logger
	ErrLogger  Logger
}

func (p *PkgDirDaemon) init() error {
	if p.DbName == "" {
		return errors.New("pkg dir daemon db name is not specified")
	}
	if p.WatchDir == "" {
		return errors.New("pkg dir daemon watch dir is not specified")
	}
	if p.MkDirMode == 0 {
		p.MkDirMode = os.ModePerm
	}
	if p.InfoLogger == nil {
		p.InfoLogger = log.Default()
	}
	if p.ErrLogger == nil {
		p.ErrLogger = log.Default()
	}
	return nil
}

// This function is launching watcher for pacman cache directory, and constatly
// adding new arch packages to database in watched directory.
func (p PkgDirDaemon) Run() error {
	err := p.init()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(p.WatchDir, os.ModePerm); err != nil {
		return err
	}

	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move)
	// TODO add recursive watcher.
	if err := w.Add(p.WatchDir); err != nil {
		return err
	}

	go func() {
		for event := range w.Event {
			file := event.FileInfo.Name()
			if strings.HasSuffix(file, pkgext) {
				var b bytes.Buffer
				err := pacman.RepoAdd(
					path.Join(p.WatchDir, p.DbName+dbext),
					path.Join(p.WatchDir, file),
					pacman.RepoAddOptions{
						Sudo:             true,
						New:              true,
						PreventDowngrade: true,
						Stdout:           &b,
						Stderr:           &b,
					},
				)
				if err != nil {
					p.ErrLogger.Printf(
						"unable to add package %s to %s in %s, err: %s",
						file, p.DbName, p.WatchDir, b.String(),
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
