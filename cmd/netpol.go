package cmd

import (
	"github.com/camilocot/kube-ns/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// netpolConfigCmd represents the netpol subcommand
var netpolConfigCmd = &cobra.Command{
	Use:   "netpol",
	Short: "netpol configuration",
	Long:  `netpol configuration via annotations`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		enabled, err := cmd.Flags().GetBool("enabled")
		if err == nil {
			conf.NetPol.Enabled = enabled

		} else {
			logrus.Fatal(err)
		}

		annotation, err := cmd.Flags().GetString("annotation")
		if err == nil {
			conf.NetPol.Annotation = annotation

		} else {
			logrus.Fatal(err)
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	netpolConfigCmd.Flags().Bool("enabled", true, "Enable or disable netpol creation)")
	netpolConfigCmd.Flags().String("annotation", "kubens/netpol.recipe", "Namespace annotation setting the netpol recipe (deny-all)")
}
