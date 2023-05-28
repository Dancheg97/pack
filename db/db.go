// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package db

type Database interface {
	List() ([]string, error)
	Validate(name string, password string) bool
	Add(name string, password string) error
}
