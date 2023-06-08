// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

// Source for public keys, which can be used for signature verification.
type PubkeySource interface {
	Get(email string) ([]string, error)
}
