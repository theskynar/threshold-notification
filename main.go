package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/theskynar/thresold-notification/envs"
	"github.com/theskynar/thresold-notification/infra/db"
	"github.com/theskynar/thresold-notification/infra/glog"
)

var database *db.PG

func init() {
	glog.SetupLog()

	log.Info("Checking the environment variables...")
	if err := envs.SetupEnv(); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to setup the envs"))
	}

	log.Info("Starting the application")

}

func main() {
	database = &db.PG{}
	if err := database.Open(); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to open the connection with the database"))
	}

	go func() {
		for {
			select {
			case err := <-database.ErrCh:
				log.Fatal(err.Error)
			case info := <-database.InfoCh:
				log.Info(info)
			}
		}
	}()

	// Listening the OS exit signals, to apply the shutdown the application gracefully
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGTERM)
	signal.Notify(shutdownCh, syscall.SIGHUP)
	signal.Notify(shutdownCh, syscall.SIGINT)

	go func() {
		select {
		case signal := <-shutdownCh:
			log.Infof("Got OS Signal %s", signal)
			log.Info("Shutting down...")

			if err := database.Close(); err != nil {
				log.Error(err)
				os.Exit(1)
				return
			}

			log.Info("Shutdown successfully")

			os.Exit(0)
		}
	}()

	log.Info("Started the application")

	for {
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println("Hello World!")
	}
}
