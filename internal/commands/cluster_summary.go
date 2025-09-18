package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service/cluster"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/txt/report"
)

// ClusterSummaryCommand prints a minimal cluster summary (Portal-only).
var ClusterSummaryCommand = &cli.Command{
	Name:   "summary",
	Usage:  "Shows cluster summary (Portal-only)",
	Flags:  append(report.CliFlags, JsonFlag),
	Action: clusterSummaryAction,
}

func clusterSummaryAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		if !conf.IsPortal() {
			return fmt.Errorf("cluster summary is only available on a Portal node")
		}

		r, err := reg.NewFileRegistry(conf)
		if err != nil {
			return err
		}

		nodes, _ := r.List()

		resp := cluster.SummaryResponse{
			PortalUUID: conf.PortalUUID(),
			Nodes:      len(nodes),
			DB:         cluster.DBInfo{Driver: conf.DatabaseDriverName(), Host: conf.DatabaseHost(), Port: conf.DatabasePort()},
			Time:       time.Now().UTC().Format(time.RFC3339),
		}

		if ctx.Bool("json") {
			b, _ := json.Marshal(resp)
			fmt.Println(string(b))
			return nil
		}

		cols := []string{"Portal UUID", "Nodes", "DB Driver", "DB Host", "DB Port", "Time"}
		rows := [][]string{{resp.PortalUUID, fmt.Sprintf("%d", resp.Nodes), resp.DB.Driver, resp.DB.Host, fmt.Sprintf("%d", resp.DB.Port), resp.Time}}
		out, err := report.RenderFormat(rows, cols, report.CliFormat(ctx))
		fmt.Printf("\n%s\n", out)
		return err
	})
}
