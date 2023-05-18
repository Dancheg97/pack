// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package cmd

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/prnt"
	"fmnx.su/core/pack/tmpl"
	"github.com/spf13/cobra"
)

// This package contains all CLI commands that might be executed by user.
// Each file contains a single command, including root cmd.

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   tmpl.SearchShort,
	Long:    tmpl.SearchLong,
	Run:     Search,
}

// Search for packages in connected search sources.
func Search(cmd *cobra.Command, args []string) {
	for _, v := range args {
		prnt.Custom([]prnt.ColoredMessage{
			{
				Message: "> Search request",
				Color:   prnt.COLOR_GREEN,
			},
			{
				Message: " - ",
				Color:   prnt.COLOR_WHITE,
			},
			{
				Message: v,
				Color:   prnt.COLOR_CYAN,
			},
		})
		for _, ss := range config.SearchSources {
			prnt.Custom([]prnt.ColoredMessage{
				{
					Message: " => Seach in",
					Color:   prnt.COLOR_BLUE,
				},
				{
					Message: " - ",
					Color:   prnt.COLOR_WHITE,
				},
				{
					Message: ss.Name,
					Color:   prnt.COLOR_YELLOW,
				},
			})
			rez, err := UrlSearch(v, ss.Url, ss.Field)
			CheckErr(err)
			for i, pkg := range rez {
				fmt.Printf("  %d - %s%s\n", i, ss.Prefix, pkg)

			}
		}
	}
}

// This package implements funcitons related to UrlSearch of arch packages in
// different compatible formats. Safe for concurrent usage.
func UrlSearch(req string, url string, field string) ([]string, error) {
	resp, err := http.Get(strings.ReplaceAll(url, `{{package}}`, req))
	if err != nil {
		return nil, err
	}
	bodystr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ejectSearchFields(field, string(bodystr)), nil
}

// Eject values from json by field.
func ejectSearchFields(field string, json string) []string {
	splt := strings.Split(json, "\""+field+"\"")
	var rez []string
	for i, v := range splt {
		if i == 0 {
			continue
		}
		vsplt := strings.Split(v, "\"")
		rez = append(rez, vsplt[1])
	}
	return rez
}
