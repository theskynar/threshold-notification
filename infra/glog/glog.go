package glog

import (
	"os"

	"github.com/theskynar/thresold-notification/envs"

	log "github.com/sirupsen/logrus"
)

func SetupLog() {
	if envs.Variables.Name == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}
	log.SetOutput(os.Stdout)
}
