package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const kubensConfigFile = ".kubens.yaml"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "modify kubens configuration",
	Long: `
config command allows configuration of ~/.kubens.yaml for running kubens`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add config to ~/.kubens.yaml",
	Long: `
Adds config to ~/.kubens.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "view ~/.kubens.yaml",
	Long: `
Display the contents of the contents of ~/.kubens.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Contents of ~/.kubens.yaml")
		configFile, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), kubensConfigFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(configFile))
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configAddCmd,
		configViewCmd,
	)

	configAddCmd.AddCommand(
		netpolConfigCmd,
	)
}
