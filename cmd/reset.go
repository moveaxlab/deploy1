package cmd

import (
	"golang.org/x/sync/errgroup"

	"github.com/moveaxlab/deploy1/argo"
	"github.com/moveaxlab/deploy1/config"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset (SERVICE ...)",
	Short: "Remove image tag overrides for one or more services",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.AutocompleteService(args, toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		applyDebugFlag(cmd)

		services, err := config.ValidateServiceName(args)
		checkNoError(err)

		resetConfig, err := getResetFlags(cmd)
		checkNoError(err)

		var g errgroup.Group
		for _, service := range services {
			service := service
			g.Go(func() error {
				return argo.Reset(config.GetServiceName(service, resetConfig.env), resetConfig.env, config.GetImageTagParameter(service), resetConfig.wait)
			})
		}

		checkNoError(g.Wait())
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	addDebugFlag(resetCmd)
	addResetFlags(resetCmd)
}
