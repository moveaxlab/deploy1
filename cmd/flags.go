package cmd

import (
	"fmt"
	"github.com/moveaxlab/deploy1/config"
	"github.com/moveaxlab/deploy1/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	debugFlag       = "debug"
	tagFlag         = "tag"
	envFlag         = "env"
	buildArgsFlag   = "build-args"
	deployFlag      = "deploy"
	noBundleFlag    = "no-bundle"
	noCacheFlag     = "no-cache"
	noImageTagCheck = "no-image-tag-check"
)

func addDebugFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(debugFlag, false, "print debug messages")
}

func applyDebugFlag(cmd *cobra.Command) {
	debugEnabled, _ := cmd.Flags().GetBool(debugFlag)

	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}
}

func addBaseFlags(cmd *cobra.Command) {
	defaultTag, _ := tag.GetTag("")

	cmd.Flags().String(
		tagFlag,
		defaultTag,
		"docker image tag",
	)

	cmd.Flags().String(
		envFlag,
		string(config.Config.DefaultEnvironment),
		fmt.Sprintf("environment for build/deploy"),
	)

	_ = cmd.RegisterFlagCompletionFunc(envFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.AutocompleteEnvironment(toComplete), cobra.ShellCompDirectiveDefault
	})
}

type baseFlags struct {
	tag string
	env config.Environment
}

func getBaseFlags(cmd *cobra.Command) (*baseFlags, error) {
	tag, err := cmd.Flags().GetString(tagFlag)

	if err != nil {
		return nil, fmt.Errorf("invalid tag: %w", err)
	}

	rawEnv, err := cmd.Flags().GetString(envFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid env: %w", err)
	}

	env, err := config.ValidateEnvironment(rawEnv)
	if err != nil {
		return nil, fmt.Errorf("invalid env: %w", err)
	}

	return &baseFlags{
		tag: tag,
		env: env,
	}, nil
}

func addBuildFlags(cmd *cobra.Command) {
	cmd.Flags().String(
		buildArgsFlag,
		"",
		"custom build arguments for docker build",
	)

	cmd.Flags().Bool(
		deployFlag,
		false,
		"deploy service(s) after build",
	)

	cmd.Flags().Bool(
		noBundleFlag,
		false,
		"skip prepare/bundle/cleanup steps (assumes everything is ready for docker build)",
	)

	cmd.Flags().Bool(
		noCacheFlag,
		false,
		"don't use docker cache during build",
	)
}

type buildArgs struct {
	buildArgs string
	deploy    bool
	noBundle  bool
	noCache   bool
}

func getBuildFlags(cmd *cobra.Command) (*buildArgs, error) {
	buildArguments, err := cmd.Flags().GetString(buildArgsFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid build arguments flag: %w", err)
	}

	deploy, err := cmd.Flags().GetBool(deployFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid deploy flag: %w", err)
	}

	noBundle, err := cmd.Flags().GetBool(noBundleFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid no bundle flag: %w", err)
	}

	noCache, err := cmd.Flags().GetBool(noCacheFlag)
	if err != nil {
		return nil, fmt.Errorf("invalid no cache flag: %w", err)
	}

	return &buildArgs{
		buildArgs: buildArguments,
		deploy:    deploy,
		noBundle:  noBundle,
		noCache:   noCache,
	}, nil
}

func addDeployFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(
		noImageTagCheck,
		false,
		"skip image tag check before deployment (use at your own risk)",
	)
}

type deployFlags struct {
	noImageTagCheck bool
}

func getDeployFlags(cmd *cobra.Command) (*deployFlags, error) {
	noImageTagCheck, err := cmd.Flags().GetBool(noImageTagCheck)
	if err != nil {
		return nil, fmt.Errorf("invalid image tag check flag: %w", err)
	}

	return &deployFlags{
		noImageTagCheck: noImageTagCheck,
	}, nil
}
