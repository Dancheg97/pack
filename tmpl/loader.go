// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/ioprogress"
	"golang.org/x/term"
)

type LoaderParameters struct {
	Current int
	Total   int
	Msg     string
	Output  io.Writer
}

// Function that will give terminal drawer for provided message, that can be
// further used in different IO operations.
func Loader(p *LoaderParameters) func(int64, int64) error {
	width, _, err := term.GetSize(0)
	if err != nil {
		return nil
	}
	prefix := fmt.Sprintf("(%d/%d) %s", p.Current, p.Total, p.Msg)

	if len(prefix) > width {
		return ioprogress.DrawTerminalf(p.Output, func(i1, i2 int64) string {
			return prefix[:width-3] + "..."
		})
	}

	loaderwidth := int(float64(width) * 0.35)
	paddingWidth := width - len(prefix) - loaderwidth - 7

	if paddingWidth < 0 {
		return ioprogress.DrawTerminalf(p.Output, func(i1, i2 int64) string {
			return prefix
		})
	}

	padding := strings.Repeat(" ", paddingWidth)
	return ioprogress.DrawTerminalf(p.Output, func(progress, total int64) string {
		prcntg := float32(progress) / float32(total) * 100

		current := int((float64(progress) / float64(total)) * float64(loaderwidth))
		bar := fmt.Sprintf(
			"[%s%s]", strings.Repeat("#", int(current)),
			strings.Repeat("-", int(loaderwidth-current)),
		)

		return fmt.Sprintf("%s%s%s %.0f", prefix, padding, bar, prcntg) + "%"
	})
}
