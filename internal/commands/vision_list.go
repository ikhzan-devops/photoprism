package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize/english"
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/pkg/txt/report"
)

// VisionListCommand configures the command name, flags, and action.
var VisionListCommand = &cli.Command{
	Name:   "ls",
	Usage:  "Lists the configured computer vision models",
	Flags:  append(report.CliFlags),
	Action: visionListAction,
}

// visionListAction displays existing user accounts.
func visionListAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		var rows [][]string

		cols := []string{
			"Type",
			"Name",
			"Version",
			"Resolution",
			"Provider",
			"Service Endpoint",
			"Request Format",
			"Response Format",
			"Options",
			"Tags",
			"Disabled",
		}

		// Show log message.
		log.Infof("found %s", english.Plural(len(vision.Config.Models), "model", "models"))

		if n := len(vision.Config.Models); n == 0 {
			return nil
		} else {
			rows = make([][]string, n)
		}

		// Display report.
		for i, model := range vision.Config.Models {
			modelUri, modelMethod := model.Endpoint()
			tags := ""

			_, name, version := model.Model()

			if model.TensorFlow != nil && model.TensorFlow.Tags != nil {
				tags = strings.Join(model.TensorFlow.Tags, ", ")
			}

			if model.Default {
				version = "default"
			}

			var options []byte
			if o := model.GetOptions(); o != nil {
				options, _ = json.Marshal(*o)
			}

			var responseFormat, requestFormat string

			if modelUri != "" && modelMethod != "" {
				if f := strings.TrimSpace(string(model.EndpointRequestFormat())); f != "" {
					requestFormat = f
				}

				if f := strings.TrimSpace(string(model.EndpointResponseFormat())); f != "" {
					responseFormat = f
				}
			}

			provider := model.ProviderName()

			rows[i] = []string{
				model.Type,
				name,
				version,
				fmt.Sprintf("%d", model.Resolution),
				provider,
				fmt.Sprintf("%s %s", modelMethod, modelUri),
				requestFormat,
				responseFormat,
				string(options),
				tags,
				report.Bool(model.Disabled, report.Yes, report.No),
			}
		}

		result, err := report.RenderFormat(rows, cols, report.CliFormat(ctx))

		fmt.Printf("\n%s\n", result)

		return err
	})
}
