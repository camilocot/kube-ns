package client

import (
	"log"

	"github.com/camilocot/kube-ns/config"
	"github.com/camilocot/kube-ns/pkg/controller"
	"github.com/camilocot/kube-ns/pkg/handlers"
)

// Run runs the event loop processing with given handler
func Run(conf *config.Config) {

	var eventHandler = ParseEventHandler(conf)
	controller.Start(conf, eventHandler)
}

// ParseEventHandler returns the respective handler object specified in the config file.
func ParseEventHandler(conf *config.Config) handlers.Handler {

	var eventHandler handlers.Handler

	eventHandler = new(handlers.HandlerConfig)
	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}
	return eventHandler
}
