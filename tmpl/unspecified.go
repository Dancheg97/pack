// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package tmpl

// This package contains string templates for utility command output.
// Output might be modified based on configuration - prettyprint.

const DbAdd = `[%s]
SigLevel = Optional TrustAll
Server = https://%s/%s
`
