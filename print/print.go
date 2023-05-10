// Copyright 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package print

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
	WHITE  Color = iota
	RED    Color = iota
	BLUE   Color = iota
	GREEN  Color = iota
	YELLOW Color = iota
)

type ColoredMessage struct {
	Message string
	Color   Color
}

// This function will concatenate all provided messages and print them as
// single string.
func Custom(msgs []ColoredMessage) {
	var rez string
	for _, msg := range msgs {
		switch msg.Color {
		case WHITE:
			rez += msg.Message
		case RED:
			rez += color.RedString(msg.Message)
		case BLUE:
			rez += color.BlueString(msg.Message)
		case GREEN:
			rez += color.GreenString(msg.Message)
		case YELLOW:
			rez += color.YellowString(msg.Message)
		}
	}
	fmt.Println(rez)
}
