package commands

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/txt/report"
)

// ClusterNodesShowCommand shows node details.
var ClusterNodesShowCommand = &cli.Command{
	Name:      "show",
	Usage:     "Shows node details (Portal-only)",
	ArgsUsage: "<id|name>",
	Flags:     append(report.CliFlags, JsonFlag),
	Action:    clusterNodesShowAction,
}

func clusterNodesShowAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		if !conf.IsPortal() {
			return fmt.Errorf("node show is only available on a Portal node")
		}

		key := ctx.Args().First()
		if key == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		r, err := reg.NewFileRegistry(conf)
		if err != nil {
			return err
		}

		// Resolve by id first, then by normalized name.
		n, getErr := r.Get(key)
		if getErr != nil {
			name := clean.TypeLowerDash(key)
			if name == "" {
				return fmt.Errorf("invalid node identifier")
			}
			n, getErr = r.FindByName(name)
		}
		if getErr != nil || n == nil {
			return fmt.Errorf("node not found")
		}

		opts := reg.NodeOpts{IncludeInternalURL: true, IncludeDBMeta: true}
		dto := reg.BuildClusterNode(*n, opts)

		if ctx.Bool("json") {
			b, _ := json.Marshal(dto)
			fmt.Println(string(b))
			return nil
		}

		cols := []string{"ID", "Name", "Type", "Internal URL", "DB Name", "DB User", "DB Last Rotated", "Created At", "Updated At"}
		var dbName, dbUser, dbRot string
		if dto.DB != nil {
			dbName, dbUser, dbRot = dto.DB.Name, dto.DB.User, dto.DB.DBLastRotatedAt
		}
		rows := [][]string{{dto.ID, dto.Name, dto.Type, dto.InternalURL, dbName, dbUser, dbRot, dto.CreatedAt, dto.UpdatedAt}}
		out, err := report.RenderFormat(rows, cols, report.CliFormat(ctx))
		fmt.Printf("\n%s\n", out)
		return err
	})
}
