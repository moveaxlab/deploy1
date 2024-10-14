package cmd

import (
	"github.com/moveaxlab/deploy1/argo"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/docker"
	"github.com/moveaxlab/deploy1/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy (SERVICE ...)",
	Short: "Deploy one or more services",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.AutocompleteService(args, toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		applyDebugFlag(cmd)

		services, err := config.ValidateServiceName(args)
		checkNoError(err)

		baseConfig, err := getBaseFlags(cmd)
		checkNoError(err)

		deployFlags, err := getDeployFlags(cmd)
		checkNoError(err)

		actualTag, err := tag.GetTag(baseConfig.tag)
		checkNoError(err)

		for _, service := range services {
			if !deployFlags.noImageTagCheck {
				tagExists, err := docker.TagExists(service, baseConfig.env, actualTag)
				checkNoError(err)
				if !tagExists {
					log.Infof("skipping deployment of service %s: tag %s does not exist", service, actualTag)
					continue
				}
			}
			err = argo.Deploy(config.GetServiceName(service, baseConfig.env), actualTag, baseConfig.env, config.GetImageTagParameter(service))
			checkNoError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	addDebugFlag(deployCmd)
	addBaseFlags(deployCmd)
	addDeployFlags(deployCmd)
}
