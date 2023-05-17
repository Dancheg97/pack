// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package search

import (
	"io"
	"net/http"
	"strings"
)

// This package implements funcitons related to search of arch packages in
// different compatible formats. Safe for concurrent usage.
func Search(req string, url string, field string) ([]string, error) {
	resp, err := http.Get(strings.ReplaceAll(url, `{{package}}`, req))
	if err != nil {
		return nil, err
	}
	bodystr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ejectFields(field, string(bodystr)), nil
}

// Eject values from json by field.
func ejectFields(field string, json string) []string {
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
