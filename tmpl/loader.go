// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"os"
	"path"
	"strings"

	bar "github.com/mitchellh/ioprogress"
	"golang.org/x/term"
)

func Loader(registry, owner, pkg string, i, t int) func(int64, int64) error {
	w, _, err := term.GetSize(0)
	if err != nil {
		return nil
	}
	w = w - 24 - len(registry) - len(owner) - len(pkg)
	msg := fmt.Sprintf("(%d/%d) Package %s to %s", i, t, pkg, path.Join(registry, owner))
	return bar.DrawTerminalf(os.Stdout, func(progress, total int64) string {
		percentage := float32(progress) / float32(total)
		fill := strings.Repeat("#", int(float32(w)*0.52*percentage))
		rest := strings.Repeat("-", int(float32(w)*0.52*(1-percentage)))
		prg := percentage * 100
		return fmt.Sprintf("%s %*s %.0f", msg, w, "["+fill+rest+"]", prg) + "%"
	})
}
