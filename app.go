package app

import (
	log "github.com/sirupsen/logrus"
	"io"
)

type App struct {
}

func (a *App) Run(writer io.Writer) error {
	log.Info("Hello, world!")
	return nil
}
