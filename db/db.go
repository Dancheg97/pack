// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package db

type Database interface {
	List() ([]User, error)
	Update(u User) error
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
