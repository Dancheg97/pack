// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package tmpl

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	bar "github.com/mitchellh/ioprogress"
)

func Loader(registry string, owner string, pkg string) func(int64, int64) error {
	pre := Dots + color.New(color.Bold).Sprintf(" Pushing package:  ")
	done := Dots + color.New(color.Bold).Sprintf(" Package uploaded: ")
	pkg = path.Join(registry, owner, pkg)
	return bar.DrawTerminalf(os.Stdout, func(progress, total int64) string {
		prg := float32(progress) / float32(total)
		switch {
		case prg > 0.9999:
			return fmt.Sprintf("%s[=======================================] %s", done, pkg)
		case prg > 0.975:
			return fmt.Sprintf("%s[======================================>] %s", pre, pkg)
		case prg > 0.95:
			return fmt.Sprintf("%s[=====================================>-] %s", pre, pkg)
		case prg > 0.925:
			return fmt.Sprintf("%s[====================================>--] %s", pre, pkg)
		case prg > 0.9:
			return fmt.Sprintf("%s[===================================>---] %s", pre, pkg)
		case prg > 0.875:
			return fmt.Sprintf("%s[==================================>----] %s", pre, pkg)
		case prg > 0.85:
			return fmt.Sprintf("%s[=================================>-----] %s", pre, pkg)
		case prg > 0.825:
			return fmt.Sprintf("%s[================================>------] %s", pre, pkg)
		case prg > 0.8:
			return fmt.Sprintf("%s[===============================>-------] %s", pre, pkg)
		case prg > 0.775:
			return fmt.Sprintf("%s[==============================>--------] %s", pre, pkg)
		case prg > 0.75:
			return fmt.Sprintf("%s[=============================>---------] %s", pre, pkg)
		case prg > 0.725:
			return fmt.Sprintf("%s[============================>----------] %s", pre, pkg)
		case prg > 0.7:
			return fmt.Sprintf("%s[===========================>-----------] %s", pre, pkg)
		case prg > 0.675:
			return fmt.Sprintf("%s[==========================>------------] %s", pre, pkg)
		case prg > 0.65:
			return fmt.Sprintf("%s[=========================>-------------] %s", pre, pkg)
		case prg > 0.625:
			return fmt.Sprintf("%s[========================>--------------] %s", pre, pkg)
		case prg > 0.6:
			return fmt.Sprintf("%s[=======================>---------------] %s", pre, pkg)
		case prg > 0.575:
			return fmt.Sprintf("%s[======================>----------------] %s", pre, pkg)
		case prg > 0.55:
			return fmt.Sprintf("%s[=====================>-----------------] %s", pre, pkg)
		case prg > 0.525:
			return fmt.Sprintf("%s[====================>------------------] %s", pre, pkg)
		case prg > 0.5:
			return fmt.Sprintf("%s[===================>-------------------] %s", pre, pkg)
		case prg > 0.475:
			return fmt.Sprintf("%s[==================>--------------------] %s", pre, pkg)
		case prg > 0.45:
			return fmt.Sprintf("%s[=================>---------------------] %s", pre, pkg)
		case prg > 0.425:
			return fmt.Sprintf("%s[================>----------------------] %s", pre, pkg)
		case prg > 0.4:
			return fmt.Sprintf("%s[===============>-----------------------] %s", pre, pkg)
		case prg > 0.375:
			return fmt.Sprintf("%s[==============>------------------------] %s", pre, pkg)
		case prg > 0.35:
			return fmt.Sprintf("%s[=============>-------------------------] %s", pre, pkg)
		case prg > 0.325:
			return fmt.Sprintf("%s[============>--------------------------] %s", pre, pkg)
		case prg > 0.3:
			return fmt.Sprintf("%s[===========>---------------------------] %s", pre, pkg)
		case prg > 0.275:
			return fmt.Sprintf("%s[==========>----------------------------] %s", pre, pkg)
		case prg > 0.25:
			return fmt.Sprintf("%s[=========>-----------------------------] %s", pre, pkg)
		case prg > 0.225:
			return fmt.Sprintf("%s[========>------------------------------] %s", pre, pkg)
		case prg > 0.2:
			return fmt.Sprintf("%s[=======>-------------------------------] %s", pre, pkg)
		case prg > 0.175:
			return fmt.Sprintf("%s[======>--------------------------------] %s", pre, pkg)
		case prg > 0.15:
			return fmt.Sprintf("%s[=====>---------------------------------] %s", pre, pkg)
		case prg > 0.125:
			return fmt.Sprintf("%s[====>----------------------------------] %s", pre, pkg)
		case prg > 0.1:
			return fmt.Sprintf("%s[===>-----------------------------------] %s", pre, pkg)
		case prg > 0.075:
			return fmt.Sprintf("%s[==>------------------------------------] %s", pre, pkg)
		case prg > 0.05:
			return fmt.Sprintf("%s[=>-------------------------------------] %s", pre, pkg)
		case prg > 0.025:
			return fmt.Sprintf("%s[>--------------------------------------] %s", pre, pkg)
		default:
			return fmt.Sprintf("%s[---------------------------------------] %s", pre, pkg)
		}
	})
}
