// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package main

// Package description.
//
// cmd - commands, each file represents a single command
//
// config - CLI configuration
//
// git - library for accessing git functionality
//
// pacman - library for accessing pacman and makepkg functionality
//
// pack - all operations related to pack database and outputs
//
// prnt - utility for pretty printing
//
// system - library for executing system calls and file operations
//
// tmpl - string templates

import "fmnx.io/core/pack/cmd"

func main() {
	cmd.Execute()
}
