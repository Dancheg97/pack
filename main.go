// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package main

import (
	"fmt"
	"os"

	"fmnx.su/core/pack/cmd"
)

func main() {
	err := cmd.Command.Execute()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
