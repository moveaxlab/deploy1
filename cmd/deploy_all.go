package cmd

import (
	"github.com/moveaxlab/deploy1/argo"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/docker"
	"github.com/moveaxlab/deploy1/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deployAllCmd = &cobra.Command{
	Use:   "deploy-all",
	Short: "Deploy all services",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		applyDebugFlag(cmd)

		baseConfig, err := getBaseFlags(cmd)
		checkNoError(err)

		actualTag, err := tag.GetTag(baseConfig.tag)
		checkNoError(err)

		deployFlags, err := getDeployFlags(cmd)
		checkNoError(err)

		for _, service := range config.GetAllServices() {
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
	rootCmd.AddCommand(deployAllCmd)
	addDebugFlag(deployAllCmd)
	addBaseFlags(deployAllCmd)
	addDeployFlags(deployAllCmd)
}
