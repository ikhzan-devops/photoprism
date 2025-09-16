package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/txt/report"
)

// ClusterNodesCommands groups node subcommands.
var ClusterNodesCommands = &cli.Command{
	Name:  "nodes",
	Usage: "Node registry subcommands",
	Subcommands: []*cli.Command{
		ClusterNodesListCommand,
		ClusterNodesShowCommand,
		ClusterNodesModCommand,
		ClusterNodesRemoveCommand,
		ClusterNodesRotateCommand,
	},
}

// ClusterNodesListCommand lists registered nodes.
var ClusterNodesListCommand = &cli.Command{
	Name:      "ls",
	Usage:     "Lists registered cluster nodes (Portal-only)",
	Flags:     append(append(report.CliFlags, JsonFlag), CountFlag, OffsetFlag),
	ArgsUsage: "",
	Action:    clusterNodesListAction,
}

func clusterNodesListAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		if !conf.IsPortal() {
			return fmt.Errorf("node listing is only available on a Portal node")
		}

		r, err := reg.NewFileRegistry(conf)
		if err != nil {
			return err
		}

		items, err := r.List()
		if err != nil {
			return err
		}

		// Pagination identical to API defaults.
		count := int(ctx.Uint("count"))
		if count <= 0 || count > 1000 {
			count = 100
		}
		offset := ctx.Int("offset")
		if offset < 0 {
			offset = 0
		}
		if offset > len(items) {
			offset = len(items)
		}
		end := offset + count
		if end > len(items) {
			end = len(items)
		}
		page := items[offset:end]

		// Build admin view (include internal URL and DB meta).
		opts := reg.NodeOpts{IncludeInternalURL: true, IncludeDBMeta: true}
		out := reg.BuildClusterNodes(page, opts)

		if ctx.Bool("json") {
			b, _ := json.Marshal(out)
			fmt.Println(string(b))
			return nil
		}

		cols := []string{"ID", "Name", "Type", "Labels", "Internal URL", "DB Name", "DB User", "DB Last Rotated", "Created At", "Updated At"}
		rows := make([][]string, 0, len(out))
		for _, n := range out {
			var dbName, dbUser, dbRot string
			if n.DB != nil {
				dbName, dbUser, dbRot = n.DB.Name, n.DB.User, n.DB.DBLastRotatedAt
			}
			rows = append(rows, []string{
				n.ID, n.Name, n.Type, formatLabels(n.Labels), n.InternalURL, dbName, dbUser, dbRot, n.CreatedAt, n.UpdatedAt,
			})
		}

		if len(rows) == 0 {
			log.Warnf("no nodes registered")
			return nil
		}

		result, err := report.RenderFormat(rows, cols, report.CliFormat(ctx))
		fmt.Printf("\n%s\n", result)
		return err
	})
}

func formatLabels(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ", ")
}
