package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service/cluster/provisioner"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/clean"
)

// ClusterNodesRemoveCommand deletes a node from the registry.
var ClusterNodesRemoveCommand = &cli.Command{
	Name:      "rm",
	Usage:     "Deletes a node from the registry",
	ArgsUsage: "<id|name>",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Usage: "runs the command non-interactively"},
		&cli.BoolFlag{Name: "all-ids", Usage: "delete all records that share the same UUID (admin cleanup)"},
		&cli.BoolFlag{Name: "drop-db", Aliases: []string{"d"}, Usage: "drop the node’s provisioned database and user after registry deletion"},
	},
	Hidden: true, // Required for cluster-management only.
	Action: clusterNodesRemoveAction,
}

func clusterNodesRemoveAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		if !conf.IsPortal() {
			return cli.Exit(fmt.Errorf("node delete is only available on a Portal node"), 2)
		}

		key := ctx.Args().First()

		if key == "" {
			return cli.Exit(fmt.Errorf("node id or name is required"), 2)
		}

		r, err := reg.NewClientRegistryWithConfig(conf)

		if err != nil {
			return cli.Exit(err, 1)
		}

		// Resolve to id for deletion, but also support name.
		// Resolve UUID to delete: accept uuid → clientId → name.
		var node *reg.Node

		if n, findErr := r.FindByNodeUUID(key); findErr == nil && n != nil {
			node = n
		} else if n, findErr = r.FindByClientID(key); findErr == nil && n != nil {
			node = n
		} else if name := clean.DNSLabel(key); name != "" {
			if n, findErr = r.FindByName(name); findErr == nil && n != nil {
				node = n
			}
		}

		if node == nil {
			return cli.Exit(fmt.Errorf("node not found"), 3)
		}

		uuid := node.UUID

		confirmed := RunNonInteractively(ctx.Bool("yes"))
		if !confirmed {
			prompt := promptui.Prompt{Label: fmt.Sprintf("Delete node %s?", clean.Log(uuid)), IsConfirm: true}
			if _, err := prompt.Run(); err != nil {
				log.Infof("node %s was not deleted", clean.Log(uuid))
				return nil
			}
		}

		dropDB := ctx.Bool("drop-db")
		dbName, dbUser := "", ""
		if dropDB && node.Database != nil {
			dbName = node.Database.Name
			dbUser = node.Database.User
		}

		if ctx.Bool("all-ids") {
			if err := r.DeleteAllByUUID(uuid); err != nil {
				return cli.Exit(err, 1)
			}
		} else if err := r.Delete(uuid); err != nil {
			return cli.Exit(err, 1)
		}

		if dropDB {
			if dbName == "" && dbUser == "" {
				log.Infof("node %s has been deleted (no database credentials recorded)", clean.Log(uuid))
			} else {
				dropCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := provisioner.DropCredentials(dropCtx, dbName, dbUser); err != nil {
					return cli.Exit(fmt.Errorf("failed to drop database credentials for node %s: %w", clean.Log(uuid), err), 1)
				}
				log.Infof("node %s database %s and user %s have been dropped", clean.Log(uuid), clean.Log(dbName), clean.Log(dbUser))
			}
		}

		log.Infof("node %s has been deleted", clean.Log(uuid))

		return nil
	})
}
