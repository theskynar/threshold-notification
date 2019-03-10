package envs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	GoEnv  = "GO_ENV"
	PGHost = "PG_HOST"
	PGUser = "PG_USER"
	PGPass = "PG_PASS"
	PGName = "PG_NAME"
	PGPort = "PG_PORT"
)

type Environment struct {
	Name   string
	PGHost string
	PGUser string
	PGPass string
	PGName string
	PGPort int
}

var Variables *Environment = &Environment{}

func SetupEnv() error {
	goEnv := os.Getenv(GoEnv)
	if goEnv == "" || (goEnv != "production" && goEnv != "development") {
		return errors.New(fmt.Sprintf("GO_ENV has to be either production or development, got '%s'", goEnv))
	}

	pgHost := os.Getenv(PGHost)
	if pgHost == "" {
		return errors.New("PG_HOST is required")
	}

	pgUser := os.Getenv(PGUser)
	if pgHost == "" {
		return errors.New("PG_USER is required")
	}

	pgPass := os.Getenv(PGPass)
	if pgHost == "" {
		return errors.New("PG_PASS is required")
	}

	pgName := os.Getenv(PGName)
	if pgHost == "" {
		return errors.New("PG_NAME is required")
	}

	pgPort := os.Getenv(PGPort)
	if pgPort == "" {
		Variables.PGPort = 5432
	} else {
		i, err := strconv.Atoi(pgPort)
		if err != nil {
			return errors.Wrapf(err, "Failed to convert PG_PORT with the value '%s' to integer", pgPort)
		}

		Variables.PGPort = i
	}

	Variables.Name = goEnv
	Variables.PGHost = pgHost
	Variables.PGName = pgName
	Variables.PGPass = pgPass
	Variables.PGUser = pgUser

	return nil
}
