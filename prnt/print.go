// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package prnt

// This package contains operations required for pretty output based on
// configuration.

import (
	"fmt"

	"fmnx.io/core/pack/config"
	"github.com/fatih/color"
)

// Print white message and red postfix.
func Red(white string, red string) {
	if config.DisablePrettyPrint {
		fmt.Printf(white + red + "\n")
		return
	}
	fmt.Printf(white + color.RedString(red) + "\n")
}

// Print white message and blue postfix.
func Blue(white string, blue string) {
	if config.DisablePrettyPrint {
		fmt.Printf(white + blue + "\n")
		return
	}
	fmt.Printf(white + color.BlueString(blue) + "\n")
}

// Print white message and green postfix.
func Green(white string, green string) {
	if config.DisablePrettyPrint {
		fmt.Printf(white + green + "\n")
		return
	}
	fmt.Printf(white + color.GreenString(green) + "\n")
}

// Print white message and green postfix.
func Yellow(white string, yellow string) {
	if config.DisablePrettyPrint {
		fmt.Printf(white + yellow + "\n")
		return
	}
	fmt.Printf(white + color.YellowString(yellow) + "\n")
}

// Enum for custom colored message.
type Color int

const (
	COLOR_WHITE  Color = iota
	COLOR_RED    Color = iota
	COLOR_BLUE   Color = iota
	COLOR_GREEN  Color = iota
	COLOR_YELLOW Color = iota
)

type ColoredMessage struct {
	Message string
	Color   Color
}

// This function will concatenate all provided messages and print them as
// single string.
func Custom(msgs []ColoredMessage) {
	var rez string
	if config.DisablePrettyPrint {
		for _, msg := range msgs {
			rez += msg.Message
		}
		fmt.Println(rez)
		return
	}
	for _, msg := range msgs {
		switch msg.Color {
		case COLOR_WHITE:
			rez += msg.Message
		case COLOR_RED:
			rez += color.RedString(msg.Message)
		case COLOR_BLUE:
			rez += color.BlueString(msg.Message)
		case COLOR_GREEN:
			rez += color.GreenString(msg.Message)
		case COLOR_YELLOW:
			rez += color.YellowString(msg.Message)
		}
	}
	fmt.Println(rez)
}
