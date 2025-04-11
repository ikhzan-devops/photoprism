package commands

import (
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
)

// VisionRunCommand configures the command name, flags, and action.
var VisionRunCommand = &cli.Command{
	Name:      "run",
	Usage:     "Runs a computer vision model",
	ArgsUsage: "[type]",
	Action:    visionRunAction,
	Hidden:    true,
}

// visionListAction displays existing user accounts.
func visionRunAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		log.Error("not implemented")
		return nil
	})
}
