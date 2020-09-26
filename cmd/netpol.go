package cmd

import (
	"github.com/camilocot/kubernetes-ns-default-netpol/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// netpolConfigCmd represents the netpol subcommand
var netpolConfigCmd = &cobra.Command{
	Use:   "netpol",
	Short: "specific netpol configuration",
	Long:  `specific netpol configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		recipe, err := cmd.Flags().GetString("recipe")
		if err == nil {
			if len(recipe) > 0 {
				conf.NetPol.Recipe = recipe
			}
		} else {
			logrus.Fatal(err)
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	netpolConfigCmd.Flags().StringP("recipe", "r", "", "Specify netpol recipe (deny-all or none)")
}
