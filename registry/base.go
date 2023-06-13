// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package registry

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"fmnx.su/core/pack/tmpl"
	"github.com/fatih/color"
)

// Parameters that are required to create new registry.
type Registry struct {
	// Standart and error outputs for request logging.
	Stdout io.Writer
	Stderr io.Writer

	// Directory, where temporary files will be stored.
	TmpDir string
	// Database name, typically should match the domain.
	Dbname string
	// File storage where resulting packeges will be stored.
	FileStorage
	// Source of GPG keys for incoming packages validation.
	KeyReader
	// Optional storage for package metadata, can be used for visual
	// representation for incoming packages.
	MetadataSaver
}

// File storage, that will be used to save/retrieve data related to specific
// packages.
type FileStorage interface {
	Get(key string) (io.Reader, error)
	Save(key string, content io.Writer) error
}

// Interface, that will be used to obtain all keys, that might be used to
// verify signatures for incoming packages.
type KeyReader interface {
	ReadKey(owner, email string) ([]string, error)
}

// Interface that can be used to save package metadata, for visual
// interpretation of saved package. Optional.
type MetadataSaver interface {
	UpdateMetadata(pkg string, owner string, md any) error
}

// Write header, log error and end request.
func (p *Registry) end(w http.ResponseWriter, status int, msg error) {
	errmsg := []byte(tmpl.Err + msg.Error())
	p.Stderr.Write(errmsg)
	w.WriteHeader(status)
	w.Write(errmsg)
}

// Write an announcement message with dots prefix and bold text to provided
// io.Writer.
func (p *Registry) Amsg(msg string) {
	dots := color.New(color.FgWhite, color.Bold, color.FgHiBlue).Sprintf(":: ")
	msg = color.New(color.Bold).Sprintf(msg)
	p.Stdout.Write([]byte(dots + msg + "...\n"))
}

// Write step message, with enumeration which should represent state of program
// execution.
func (p *Registry) Smsg(msg string, i, t int) {
	p.Stdout.Write([]byte(fmt.Sprintf("(%d/%d) %s...\n", i, t, msg)))
}

// Function that can parse field of arch package from pacman decsription,
// outputted by pacman -Qpi or pacman -Qi commands.
func parseArchMdField(full string, field string) string {
	splt := strings.Split(full, field)
	return strings.Split(splt[1], "\n")[0]
}
