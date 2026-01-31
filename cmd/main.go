package main

import (
	"fmt"
	"os"

	app "github.com/jimschubert/delete-artifacts"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

var version = "dev"

// noinspection GoUnusedGlobalVariable
var date = "1970-01-01T00:00:00Z"
var commit = "unknown"
var projectName = "delete-artifacts"

var opts struct {
	Owner          *string     `short:"o" help:"GitHub Owner/Org name" env:"GITHUB_ACTOR"`
	Repo           *string     `short:"r" help:"GitHub Repo name" env:"GITHUB_REPO"`
	RunId          *int64      `short:"i" name:"run-id" help:"The workflow run id from which to delete artifacts" optional:""`
	MinBytes       int64       `name:"min" help:"Minimum size in bytes. Artifacts greater than this size will be deleted." default:"50000000"`
	MaxBytes       *int64      `name:"max" help:"Maximum size in bytes. Artifacts less than this size will be deleted" optional:""`
	Name           string      `short:"n" help:"Artifact name to be deleted" default:""`
	Pattern        string      `short:"p" help:"Regex pattern (POSIX) for matching artifact name to be deleted" default:""`
	ActiveDuration string      `short:"a" name:"active" help:"Consider artifacts as 'active' within this time frame, and avoid deletion. Duration formatted such as 23h59m." default:""`
	LogLevel       string      `short:"l" name:"log-level" help:"Log level (trace, debug, info, warn, error, fatal, panic)" env:"LOG_LEVEL" default:"info"`
	DryRun         bool        `name:"dry-run" help:"Dry-run that does not perform deletions"`
	Version        VersionFlag `short:"v" help:"Display version information"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

func main() {
	ctx := kong.Parse(&opts,
		kong.Name(projectName),
		kong.Description("Delete GitHub Actions artifacts"),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s (%s)[%s]", version, commit, date),
		},
	)

	initLogging(opts.LogLevel)

	application, err := app.New(
		opts.Owner,
		opts.Repo,
		opts.RunId,
		opts.MinBytes,
		opts.MaxBytes,
		opts.Name,
		opts.Pattern,
		opts.ActiveDuration,
		opts.DryRun)
	ctx.FatalIfErrorf(err, "unable to construct application with specific parameters.")
	err = application.Run()
	ctx.FatalIfErrorf(err, "execution failed.")

	log.Info("Run complete.")
}

func initLogging(level string) {
	ll, err := log.ParseLevel(level)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.SetOutput(os.Stderr)
}
