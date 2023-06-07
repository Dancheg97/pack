// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

func formOptions[Options any](arr []Options, getdefault func() *Options) *Options {
	if len(arr) != 1 {
		return getdefault()
	}
	return &arr[0]
}
