package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v32/github"
)

// App is the main application container
type App struct {
	Owner    *string
	Repo     *string
	RunId    *int64
	MinBytes int64
	MaxBytes *int64
	Name     string
	Pattern  string
	DryRun   bool
	context  *context.Context
	client   *github.Client
}

// Run the application
func (a *App) Run() error {
	err := a.checkPreconditions()
	if err != nil {
		return err
	}

	executionContext, cancel := context.WithTimeout(*a.context, 2*time.Minute)
	defer cancel()

	wg := sync.WaitGroup{}
	doneChan := make(chan error)
	errorChan := make(chan error)
	itemsChan := make(chan []*github.Artifact)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func(page int) {
		a.retrieveArtifactsByPage(&wg, &executionContext, page, itemsChan, errorChan)
	}(1)

	go wait(doneChan, &wg)

	all := make([]*github.Artifact, 0)
	for {
		select {
		case sig := <-signalChannel:
			log.Warn("Received signal: ", sig)
			os.Exit(0)
		case e := <-errorChan:
			return e
		case items := <-itemsChan:
			if items != nil {
				filtered := a.filterArtifacts(items)
				if len(filtered) > 0 {
					for _, artifact := range filtered {
						all = append(all, artifact)
					}
				}
			}
		case <-doneChan:
			if len(all) == 0 {
				log.Info("No artifacts to delete!")
			} else {
				if a.DryRun {
					for _, artifact := range all {
						log.WithFields(log.Fields{"size": artifact.GetSizeInBytes(), "name": artifact.GetName()}).Warn("DryRun: would have deleted the artifact")
					}
				} else {
					// perform the deletions. Synchronously is fine here.
					for _, artifact := range all {
						log.WithFields(log.Fields{"size": artifact.GetSizeInBytes(), "name": artifact.GetName()}).Info("Deleting artifact")
						_, err := a.client.Actions.DeleteArtifact(executionContext, *a.Owner, *a.Repo, artifact.GetID())
						if err != nil {
							log.Warnf("Error deleting %s (artifact ID %d), ignoring…", artifact.GetName(), artifact.GetID())
						}
					}
				}
			}
			return nil
		}
	}
}

func (a *App) filterArtifacts(artifacts []*github.Artifact) []*github.Artifact {
	filtered := make([]*github.Artifact, 0)
	for _, artifact := range artifacts {
		log.WithFields(log.Fields{"size": artifact.GetSizeInBytes(), "name": artifact.GetName()}).Debug("Iterating artifact.")
		shouldAdd := false
		size := artifact.GetSizeInBytes()
		// note MinBytes is required. it will short-circuit all other checks
		if size >= a.MinBytes {
			shouldAdd = true
		}

		if shouldAdd && a.MaxBytes != nil && size > *a.MaxBytes {
			shouldAdd = false
		}

		if shouldAdd && len(a.Name) > 0 && artifact.GetName() == a.Name {
			shouldAdd = true
		}

		if shouldAdd && len(a.Pattern) > 0 {
			re := regexp.MustCompile(a.Pattern)
			shouldAdd = re.MatchString(artifact.GetName())
		}

		if shouldAdd {
			filtered = append(filtered, artifact)
		}
	}
	return filtered
}

func (a *App) retrieveArtifactsByPage(wg *sync.WaitGroup, parent *context.Context, page int, itemsChan chan []*github.Artifact, errChan chan error) {
	ctx, timeout := context.WithTimeout(*parent, 30*time.Second)
	defer timeout()
	defer wg.Done()

	var err error
	var list *github.ArtifactList
	opts := &github.ListOptions{PerPage: 100, Page: page}
	if a.RunId != nil {
		list, _, err = a.client.Actions.ListWorkflowRunArtifacts(ctx, *a.Owner, *a.Repo, *a.RunId, opts)
	} else {
		list, _, err = a.client.Actions.ListArtifacts(ctx, *a.Owner, *a.Repo, opts)
	}

	if err != nil {
		errChan <- err
		return
	}

	if len(list.Artifacts) > 0 {
		itemsChan <- list.Artifacts

		wg.Add(1)

		go func(p int) {
			a.retrieveArtifactsByPage(wg, parent, p, itemsChan, errChan)
		}(page + 1)
	}
}

func wait(ch chan error, wg *sync.WaitGroup) {
	wg.Wait()
	ch <- nil
}

func (a *App) checkPreconditions() error {
	if len(*a.Owner) <= 1 {
		return errors.New("owner is invalid")
	}
	if len(*a.Repo) <= 1 {
		return errors.New("repo is invalid")
	}

	return nil
}

// New creates an instance of App
func New(owner *string, repo *string, runId *int64, minBytes int64, maxBytes *int64, name string, pattern string, dryRun bool) (*App, error) {
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if !found {
		return nil, errors.New("GITHUB_TOKEN environment variable is missing")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	app := &App{
		Owner:    owner,
		Repo:     repo,
		RunId:    runId,
		MinBytes: minBytes,
		MaxBytes: maxBytes,
		Name:     name,
		Pattern:  pattern,
		DryRun:   dryRun,
		context:  &ctx,
		client:   client,
	}

	return app, nil
}
