package commands

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/internal/testextras"
	"github.com/photoprism/photoprism/pkg/capture"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")

	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)
	event.AuditLog = log

	caller := "internal/commands/commands_test.go/TestMain"
	dbc, err := testextras.AcquireDBMutex(log, caller)
	if err != nil {
		log.Error("FAIL")
		os.Exit(1)
	}
	defer testextras.UnlockDBMutex(dbc.Db())

	c := config.NewTestConfig("commands")
	get.SetConfig(c)

	// Remember to close database connection.
	defer c.CloseDb()

	// Init config and connect to database.
	InitConfig = func(ctx *cli.Context) (*config.Config, error) {
		return c, c.Init()
	}

	// Run unit tests.
	beforeTimestamp := time.Now().UTC()
	code := m.Run()
	code = testextras.ValidateDBErrors(dbc.Db(), log, beforeTimestamp, code)

	testextras.ReleaseDBMutex(dbc.Db(), log, caller, code)

	os.Exit(code)
}

// NewTestContext creates a new CLI test context with the flags and arguments provided.
func NewTestContext(args []string) *cli.Context {
	// Create new command-line test app.
	app := cli.NewApp()
	app.Name = "photoprism"
	app.Usage = "PhotoPrism®"
	app.Description = ""
	app.Version = "test"
	app.Copyright = "(c) 2018-2025 PhotoPrism UG. All rights reserved."
	app.Flags = config.Flags.Cli()
	app.Commands = PhotoPrism
	app.HelpName = app.Name
	app.CustomAppHelpTemplate = ""
	app.HideHelp = true
	app.HideHelpCommand = true
	app.Action = func(*cli.Context) error { return nil }
	app.EnableBashCompletion = false
	app.Metadata = map[string]interface{}{
		"Name":    "PhotoPrism",
		"About":   "PhotoPrism®",
		"Edition": "ce",
		"Version": "test",
	}

	// Parse command test arguments.
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	LogErr(flagSet.Parse(args))

	// Create and return new test context.
	return cli.NewContext(app, flagSet, cli.NewContext(app, flagSet, nil))
}

// RunWithTestContext executes a command with a test context and returns its output.
func RunWithTestContext(cmd *cli.Command, args []string) (output string, err error) {
	// Create test context with flags and arguments.
	ctx := NewTestContext(args)

	cmd.HideHelp = false

	// Redirect the output from cli to buffer for transfer to output for testing
	var catureOutput bytes.Buffer
	oldWriter := ctx.App.Writer
	ctx.App.Writer = &catureOutput
	// Run command with test context.
	output = capture.Output(func() {
		err = cmd.Run(ctx, args...)
	})
	ctx.App.Writer = oldWriter
	output += catureOutput.String()

	return output, err
}

// NewTestContextWithParse creates a new CLI test context with the flags and arguments provided.
func NewTestContextWithParse(appArgs []string, cmdArgs []string) *cli.Context {
	// Create new command-line test app.
	app := cli.NewApp()
	app.Name = "photoprism"
	app.Usage = "PhotoPrism®"
	app.Description = ""
	app.Version = "test"
	app.Copyright = "(c) 2018-2025 PhotoPrism UG. All rights reserved."
	app.Flags = config.Flags.Cli()
	app.Commands = PhotoPrism
	app.HelpName = app.Name
	app.CustomAppHelpTemplate = ""
	app.HideHelp = true
	app.HideHelpCommand = true
	app.Action = func(*cli.Context) error { return nil }
	app.EnableBashCompletion = false
	app.Metadata = map[string]interface{}{
		"Name":    "PhotoPrism",
		"About":   "PhotoPrism®",
		"Edition": "ce",
		"Version": "test",
	}

	// Parse photoprism command arguments.
	photoprismFlagSet := flag.NewFlagSet("photoprism", flag.ContinueOnError)
	for _, f := range app.Flags {
		f.Apply(photoprismFlagSet)
	}
	LogErr(photoprismFlagSet.Parse(appArgs[1:]))

	// Parse command test arguments.
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	LogErr(flagSet.Parse(cmdArgs))

	// Create and return new test context.
	return cli.NewContext(app, flagSet, cli.NewContext(app, photoprismFlagSet, nil))
}

func RunWithProvidedTestContext(ctx *cli.Context, cmd *cli.Command, args []string) (output string, err error) {
	// Redirect the output from cli to buffer for transfer to output for testing
	var catureOutput bytes.Buffer
	oldWriter := ctx.App.Writer
	ctx.App.Writer = &catureOutput
	// Run command with test context.
	output = capture.Output(func() {
		err = cmd.Run(ctx, args...)
	})
	ctx.App.Writer = oldWriter
	output += catureOutput.String()

	return output, err
}
