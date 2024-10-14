package cmd

import (
	"github.com/moveaxlab/deploy1/argo"
	"github.com/moveaxlab/deploy1/bundle"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/docker"
	"github.com/moveaxlab/deploy1/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildAllCmd = &cobra.Command{
	Use:   "build-all",
	Short: "Build all services",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		applyDebugFlag(cmd)

		baseConfig, err := getBaseFlags(cmd)
		checkNoError(err)

		buildConfig, err := getBuildFlags(cmd)
		checkNoError(err)

		actualTag, err := tag.GetTag(baseConfig.tag)
		checkNoError(err)

		if !buildConfig.noBundle {
			err = bundle.Prepare()
			checkNoError(err)

			defer func() {
				err = bundle.Cleanup()
				checkNoError(err)
			}()
		}

		for _, service := range config.GetAllServices() {
			if !buildConfig.noBundle {
				err = bundle.Bundle(service)
				checkNoError(err)
			}
			log.Infof("building docker image for %s...", service)
			err = docker.Build(service, baseConfig.env, actualTag, buildConfig.buildArgs, buildConfig.noCache)
			checkNoError(err)
			log.Infof("service %s built, pushing docker image...", service)
			err = docker.Push(service, baseConfig.env, actualTag)
			checkNoError(err)
			log.Infof("service %s build complete", service)
		}

		if buildConfig.deploy {
			for _, service := range config.GetAllServices() {
				err = argo.Deploy(config.GetServiceName(service, baseConfig.env), actualTag, baseConfig.env, config.GetImageTagParameter(service))
				checkNoError(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildAllCmd)
	addDebugFlag(buildAllCmd)
	addBaseFlags(buildAllCmd)
	addBuildFlags(buildAllCmd)
}
