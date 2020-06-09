package cmd

import (
	"os"

	"github.com/CircleCI-Public/circleci-cli/pipeline"
	"github.com/CircleCI-Public/circleci-cli/settings"
	"github.com/spf13/cobra"
)

func newPipelineCommand(config *settings.Config) *cobra.Command {

	command := &cobra.Command{
		Use:   "pipeline",
		Short: "Operate on pipelines",
	}

	trigger := &cobra.Command{
		Use:   "trigger",
		Short: "Run a job in a container on the local machine",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return pipeline.Trigger(os.Getenv("CIRCLE_TOKEN"))
		},
	}

	command.AddCommand(trigger)

	return command
}
