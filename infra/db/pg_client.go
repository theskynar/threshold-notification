package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/theskynar/thresold-notification/envs"
)

type Event struct {
	Name  string
	Error error
}

type PG struct {
	Conn   *sql.DB
	ErrCh  chan Event
	InfoCh chan string
}

func (p *PG) Open() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		envs.Variables.PGUser,
		envs.Variables.PGPass,
		envs.Variables.PGHost,
		envs.Variables.PGPort,
		envs.Variables.PGName,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return errors.Wrap(err, "Failed to open connection with Postgres")
	}

	if err := conn.Ping(); err != nil {
		return errors.Wrap(err, "Failed to ping the connection with Postgres")
	}

	p.ErrCh = make(chan Event)
	p.InfoCh = make(chan string)

	go p.registerConnListener(connStr)

	p.Conn = conn

	return nil
}

func getEventType(eventCode pq.ListenerEventType) string {
	switch eventCode {
	case 0:
		return "ListenerEventConnected"
	case 1:
		return "ListenerEventDisconnected"
	case 2:
		return "ListenerEventReconnected"
	case 3:
		return "ListenerEventConnectionAttemptFailed"
	default:
		return "Unknown error"
	}

}

func (p *PG) registerConnListener(connStr string) {
	pq.NewListener(connStr, 10*time.Second, time.Minute, func(event pq.ListenerEventType, err error) {
		eventType := getEventType(event)

		if err != nil {
			p.ErrCh <- Event{
				Name:  eventType,
				Error: errors.Wrap(err, "Unexpected error on PG connection"),
			}
		} else {
			p.InfoCh <- fmt.Sprintf("Received event %s", eventType)
		}
	})
}

func (p *PG) Close() error {
	if err := p.Conn.Close(); err != nil {
		return errors.Wrap(err, "Failed to close the connection with PG")
	}

	close(p.InfoCh)
	close(p.ErrCh)

	return nil
}
