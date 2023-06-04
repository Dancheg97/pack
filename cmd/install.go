// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package cmd

import (
	"fmt"
	"os/exec"

	"fmnx.su/core/pack/pacman"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	AddBoolFlag(&FlagParameters{
		Cmd:   installCmd,
		Name:  "keep",
		Short: "k",
		Desc:  "do not remove database after package installation",
	})
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "🔧 install packages",
	Long: `🔧 install packages

This command will automtically connect registries adding them to pacman.conf
and install provided packages. After installation it will remove registry 
from pacman configuration.`,
	Run: Install,
}

func Install(cmd *cobra.Command, args []string) {
	rgs, err := pacman.AddRegistries(args)
	CheckErr(err)
	err = pacman.SyncList(args)
	if !viper.GetBool("keep") || err != nil {
		cmd := fmt.Sprintf("cat <<EOF > /etc/pacman.conf\n%sEOF", *rgs)
		err = exec.Command("sudo", "bash", "-c", cmd).Run()
		CheckErr(err)
	}
	CheckErr(err)
}
