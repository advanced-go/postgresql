package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	dbClient *pgxpool.Pool
)

const (
	clientLoc = PkgPath + ":Startup"
)

var clientStartup core.MessageHandler = func(msg core.Message) {
	if isReady() {
		return
	}
	start := time.Now()
	rsc := accessResource(&msg)
	credentials := accessCredentials(&msg)
	err := clientStartup2(rsc, credentials)
	if err != nil {
		core.SendReply(msg, runtime.NewStatusError(0, clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	core.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
}

// clientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func clientStartup2(rsc StartupResource, credentials StartupCredentials) error {
	if isReady() {
		return nil
	}
	if rsc.Uri == "" {
		return errors.New("database URL is empty")
	}
	// Create connection string with credentials
	s, err := connectString(rsc.Uri, credentials)
	if err != nil {
		return err
	}
	// Create pooled client and acquire connection
	dbClient, err = pgxpool.New(context.Background(), s)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create connection pool: %v\n", err))
	}
	conn, err1 := dbClient.Acquire(context.Background())
	if err1 != nil {
		clientShutdown()
		return errors.New(fmt.Sprintf("unable to acquire connection from pool: %v\n", err1))
	}
	conn.Release()
	setReady()
	return nil
}

func clientShutdown() {
	if dbClient != nil {
		resetReady()
		dbClient.Close()
		dbClient = nil
	}
}

func connectString(url string, credentials StartupCredentials) (string, error) {
	// Username and password can be in the connect string Url
	if credentials == nil {
		return url, nil
	}
	username, password, err := credentials()
	if err != nil {
		return "", errors.New(fmt.Sprintf("error accessing credentials: %v\n", err))
	}
	return fmt.Sprintf(url, username, password), nil
}

// accessCredentials - access function for Credentials in a message
func accessCredentials(msg *core.Message) StartupCredentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(StartupCredentials); ok {
			return fn
		}
	}
	return nil
}

// accessResource - access function for a resource in a message
func accessResource(msg *core.Message) StartupResource {
	if msg == nil || msg.Content == nil {
		return StartupResource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(StartupResource); ok {
			return url
		}
	}
	return StartupResource{}
}
