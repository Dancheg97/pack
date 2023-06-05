// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package server

// Simple logger interface for embedding system into other projects.
type Logger interface {
	Printf(format string, v ...any)
}

// Source for public keys, which can be used for signature verification.
type PubkeySource interface {
	Get(email string) ([][]byte, error)
}

// If provided, packages would be added in subdirectories of provided base
// directory. This should be used when you want to split packages into groups
// in single directory.
type SubdirSource interface {
	Get(email string) (string, error)
}
