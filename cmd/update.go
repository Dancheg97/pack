package cmd

import (
	"os"

	"fmnx.io/core/pack/print"
	"fmnx.io/core/pack/system"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upd", "u"},
	Short:   "🗳️  update packages",
	Long: `🗳️  update packages

You can specify packages with versions, that you need them to update to, or
provide provide just links to get latest version from default branch.

If you don't specify any arguements, all packages will be updated.

Examples:
pack update
pack update fmnx.io/core/aist@v0.21
pack update git.xmpl.sh/pkg
`,
	Run: Update,
}

func Update(cmd *cobra.Command, pkgs []string) {
	Updating = true
	if len(pkgs) == 0 {
		FullPacmanUpdate()
		FullPackUpdate()
		return
	}

}

// Perform full pacman update.
func FullPacmanUpdate() {
	o, err := system.Call("sudo pacman --noconfirm -Syu")
	if err != nil {
		print.Red("Unable to update pacman packages: ", o)
		os.Exit(1)
	}
	print.Green("Pacman update: ", "done")
}

var Updating bool

// Perform full pack update.
func FullPackUpdate() {
	mp := ReadMapping()
	var pkgs []string
	for link, _ := range mp {
		pkgs = append(pkgs, link)
	}
	Get(nil, pkgs)
	print.Green("Pack update: ", "done")
}

// Verify if all packages are installed.
func VerifyInstalled(pkgs []string) {
	
}
