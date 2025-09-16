package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service/cluster"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/txt/report"
)

var (
	rotateDBFlag     = &cli.BoolFlag{Name: "db", Usage: "rotate DB credentials"}
	rotateSecretFlag = &cli.BoolFlag{Name: "secret", Usage: "rotate node secret"}
	rotatePortalURL  = &cli.StringFlag{Name: "portal-url", Usage: "Portal base `URL` (defaults to config)"}
	rotatePortalTok  = &cli.StringFlag{Name: "portal-token", Usage: "Portal access `TOKEN` (defaults to config)"}
)

// ClusterNodesRotateCommand triggers rotation via the register endpoint.
var ClusterNodesRotateCommand = &cli.Command{
	Name:      "rotate",
	Usage:     "Rotates a node's DB and/or secret via Portal (HTTP)",
	ArgsUsage: "<id|name>",
	Flags:     append([]cli.Flag{rotateDBFlag, rotateSecretFlag, &cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Usage: "runs the command non-interactively"}, rotatePortalURL, rotatePortalTok, JsonFlag}, report.CliFlags...),
	Action:    clusterNodesRotateAction,
}

func clusterNodesRotateAction(ctx *cli.Context) error {
	return CallWithDependencies(ctx, func(conf *config.Config) error {
		key := ctx.Args().First()
		if key == "" {
			return cli.ShowSubcommandHelp(ctx)
		}

		// Determine node name. On portal, resolve id->name via registry; otherwise treat key as name.
		name := clean.TypeLowerDash(key)
		if conf.IsPortal() {
			if r, err := reg.NewFileRegistry(conf); err == nil {
				if n, err := r.Get(key); err == nil && n != nil {
					name = n.Name
				} else if n, err := r.FindByName(clean.TypeLowerDash(key)); err == nil && n != nil {
					name = n.Name
				}
			}
		}
		if name == "" {
			return fmt.Errorf("invalid node identifier")
		}

		// Portal URL and token
		portalURL := ctx.String("portal-url")
		if portalURL == "" {
			portalURL = conf.PortalUrl()
		}
		if portalURL == "" {
			portalURL = os.Getenv(config.EnvVar("portal-url"))
		}
		if portalURL == "" {
			return fmt.Errorf("portal URL is required (use --portal-url or set portal-url)")
		}
		token := ctx.String("portal-token")
		if token == "" {
			token = conf.PortalToken()
		}
		if token == "" {
			token = os.Getenv(config.EnvVar("portal-token"))
		}
		if token == "" {
			return fmt.Errorf("portal token is required (use --portal-token or set portal-token)")
		}

		// Default: rotate DB only if no flag given (safer default)
		rotateDB := ctx.Bool("db") || (!ctx.IsSet("db") && !ctx.IsSet("secret"))
		rotateSecret := ctx.Bool("secret")

		confirmed := RunNonInteractively(ctx.Bool("yes"))
		if !confirmed {
			var what string
			switch {
			case rotateDB && rotateSecret:
				what = "DB credentials and node secret"
			case rotateDB:
				what = "DB credentials"
			case rotateSecret:
				what = "node secret"
			}
			prompt := promptui.Prompt{Label: fmt.Sprintf("Rotate %s for %s?", what, clean.LogQuote(name)), IsConfirm: true}
			if _, err := prompt.Run(); err != nil {
				log.Infof("rotation cancelled for %s", clean.LogQuote(name))
				return nil
			}
		}

		body := map[string]interface{}{
			"nodeName":     name,
			"rotate":       rotateDB,
			"rotateSecret": rotateSecret,
		}
		b, _ := json.Marshal(body)

		url := stringsTrimRightSlash(portalURL) + "/api/v1/cluster/nodes/register"
		var resp cluster.RegisterResponse
		if err := postWithBackoff(url, token, b, &resp); err != nil {
			return err
		}

		if ctx.Bool("json") {
			jb, _ := json.Marshal(resp)
			fmt.Println(string(jb))
			return nil
		}

		cols := []string{"ID", "Name", "Type", "DB Name", "DB User", "Host", "Port"}
		rows := [][]string{{resp.Node.ID, resp.Node.Name, resp.Node.Type, resp.DB.Name, resp.DB.User, resp.DB.Host, fmt.Sprintf("%d", resp.DB.Port)}}
		out, _ := report.RenderFormat(rows, cols, report.CliFormat(ctx))
		fmt.Printf("\n%s\n", out)

		if (resp.Secrets != nil && resp.Secrets.NodeSecret != "") || resp.DB.Password != "" {
			fmt.Println("PLEASE WRITE DOWN THE FOLLOWING CREDENTIALS; THEY WILL NOT BE SHOWN AGAIN:")
			if resp.Secrets != nil && resp.Secrets.NodeSecret != "" && resp.DB.Password != "" {
				fmt.Printf("\n%s\n", report.Credentials("Node Secret", resp.Secrets.NodeSecret, "DB Password", resp.DB.Password))
			} else if resp.Secrets != nil && resp.Secrets.NodeSecret != "" {
				fmt.Printf("\n%s\n", report.Credentials("Node Secret", resp.Secrets.NodeSecret, "", ""))
			} else if resp.DB.Password != "" {
				fmt.Printf("\n%s\n", report.Credentials("DB User", resp.DB.User, "DB Password", resp.DB.Password))
			}
			if resp.DB.DSN != "" {
				fmt.Printf("DSN: %s\n", resp.DB.DSN)
			}
		}
		return nil
	})
}
