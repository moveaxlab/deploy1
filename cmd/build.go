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

func checkNoError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

var buildCmd = &cobra.Command{
	Use:   "build (SERVICE ...)",
	Short: "Build one or more services",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.AutocompleteService(args, toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		applyDebugFlag(cmd)

		services, err := config.ValidateServiceName(args)
		checkNoError(err)

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

		for _, service := range services {
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
			hash, err := docker.GetHash(service, baseConfig.env, actualTag)
			checkNoError(err)
			log.Infof("hash for %s@%s is:\n%s", service, actualTag, hash)
		}

		if buildConfig.deploy {
			for _, service := range services {
				err = argo.Deploy(config.GetServiceName(service, baseConfig.env), actualTag, baseConfig.env)
				checkNoError(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	addDebugFlag(buildCmd)
	addBaseFlags(buildCmd)
	addBuildFlags(buildCmd)
}
