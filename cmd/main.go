package main

import (
	"fmt"
	"os"
	"strings"

	app "github.com/jimschubert/delete-artifacts"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var version = ""

//noinspection GoUnusedGlobalVariable
var date = ""
var commit = ""
var projectName = ""

var opts struct {
	Owner          *string `short:"o" long:"owner" description:"GitHub Owner/Org name" env:"GITHUB_ACTOR"`
	Repo           *string `short:"r" long:"repo" description:"GitHub Repo name" env:"GITHUB_REPO"`
	RunId          *int64  `short:"i" long:"run-id" description:"The workflow run id from which to delete artifacts" optional:"yes"`
	MinBytes       int64   `long:"min" description:"Minimum size in bytes. Artifacts greater than this size will be deleted." optional:"yes" default:"50000000"`
	MaxBytes       *int64  `long:"max" description:"Maximum size in bytes. Artifacts less than this size will be deleted" optional:"yes"`
	Name           string  `short:"n" long:"name" description:"Artifact name to be deleted" optional:"yes" default:""`
	Pattern        string  `short:"p" long:"pattern" description:"Regex pattern (POSIX) for matching artifact name to be deleted" optional:"yes" default:""`
	ActiveDuration string  `short:"a" long:"active" description:"Consider artifacts as 'active' within this time frame, and avoid deletion. Duration formatted such as 23h59m."`
	DryRun         bool    `long:"dry-run" description:"Dry-run that does not perform deletions"`
	Version        bool    `short:"v" long:"version" description:"Display version information"`
}

const parseArgs = flags.HelpFlag | flags.PassDoubleDash

func main() {
	parser := flags.NewParser(&opts, parseArgs)
	_, err := parser.Parse()
	if err != nil {
		flagError := err.(*flags.Error)
		if flagError.Type == flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
			return
		}

		if flagError.Type == flags.ErrUnknownFlag {
			_, _ = fmt.Fprintf(os.Stderr, "%s. Please use --help for available options.\n", strings.Replace(flagError.Message, "unknown", "Unknown", 1))
			return
		}
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing command line options: %s\n", err)
		return
	}

	if opts.Version {
		fmt.Printf("%s %s (%s)\n", projectName, version, commit)
		return
	}

	initLogging()

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
	if err != nil {
		log.WithError(err).Errorf("unable to construct application with specific parameters.")
		return
	}
	err = application.Run()
	if err != nil {
		log.WithError(err).Errorf("execution failed.")
		return
	}

	log.Info("Run complete.")
}

func initLogging() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.SetOutput(os.Stderr)
}
