package main

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var newrelicApp *newrelic.Application
var newrelicAppName string
var newrelicLicense string

func InitNewRelic() {
	if gotNewRelicAppName, exists := os.LookupEnv("NEW_RELIC_APP_NAME"); exists {
		newrelicAppName = gotNewRelicAppName
	}

	if gotNewRelicLicense, exists := os.LookupEnv("NEW_RELIC_APP_LICENSE"); exists {
		newrelicLicense = gotNewRelicLicense
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(newrelicAppName),
		newrelic.ConfigLicense(newrelicLicense),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		log.Fatalln(err)
	}

	newrelicApp = app
}
