// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"os"
	"path"

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
		prg := float32(progress) / float32(total) * 100
		switch {
		case prg > 99:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#######################################]", prg) + "%"
		case prg > 97.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#######################################]", prg) + "%"
		case prg > 95:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[######################################-]", prg) + "%"
		case prg > 92.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#####################################--]", prg) + "%"
		case prg > 90:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[####################################---]", prg) + "%"
		case prg > 87.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###################################----]", prg) + "%"
		case prg > 85:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##################################-----]", prg) + "%"
		case prg > 82.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#################################------]", prg) + "%"
		case prg > 80:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[################################-------]", prg) + "%"
		case prg > 77.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###############################--------]", prg) + "%"
		case prg > 75:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##############################---------]", prg) + "%"
		case prg > 72.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#############################----------]", prg) + "%"
		case prg > 70:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[############################-----------]", prg) + "%"
		case prg > 67.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###########################------------]", prg) + "%"
		case prg > 65:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##########################-------------]", prg) + "%"
		case prg > 62.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#########################--------------]", prg) + "%"
		case prg > 60:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[########################---------------]", prg) + "%"
		case prg > 57.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#######################----------------]", prg) + "%"
		case prg > 55:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[######################-----------------]", prg) + "%"
		case prg > 52.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#####################------------------]", prg) + "%"
		case prg > 50:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[####################-------------------]", prg) + "%"
		case prg > 47.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###################--------------------]", prg) + "%"
		case prg > 45:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##################---------------------]", prg) + "%"
		case prg > 42.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#################----------------------]", prg) + "%"
		case prg > 40:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[################-----------------------]", prg) + "%"
		case prg > 37.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###############------------------------]", prg) + "%"
		case prg > 35:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##############-------------------------]", prg) + "%"
		case prg > 32.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#############--------------------------]", prg) + "%"
		case prg > 30:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[############---------------------------]", prg) + "%"
		case prg > 27.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###########----------------------------]", prg) + "%"
		case prg > 25:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##########-----------------------------]", prg) + "%"
		case prg > 22.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#########------------------------------]", prg) + "%"
		case prg > 2:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[########-------------------------------]", prg) + "%"
		case prg > 17.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#######--------------------------------]", prg) + "%"
		case prg > 15:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[######---------------------------------]", prg) + "%"
		case prg > 12.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#####----------------------------------]", prg) + "%"
		case prg > 10:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[####-----------------------------------]", prg) + "%"
		case prg > 7.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[###------------------------------------]", prg) + "%"
		case prg > 5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[##-------------------------------------]", prg) + "%"
		case prg > 2.5:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[#--------------------------------------]", prg) + "%"
		default:
			return fmt.Sprintf("%s %*s %.0f", msg, w, "[---------------------------------------]", prg) + "%"
		}
	})
}
