// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

func formOptions[Options any](arr []Options, dv *Options) *Options {
	if len(arr) != 1 {
		return dv
	}
	return &arr[0]
}

type Logger interface {
	Printf(format string, v ...interface{})
}
