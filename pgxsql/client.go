package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	dbClient *pgxpool.Pool
)

const (
	clientLoc = PkgPath + ":Startup"
)

var clientStartup messaging.MessageHandler = func(msg messaging.Message) {
	if isReady() {
		return
	}
	start := time.Now()
	rsc := accessResource(&msg)
	credentials := accessCredentials(&msg)
	err := clientStartup2(rsc, credentials)
	if err != nil {
		messaging.SendReply(msg, runtime.NewStatusError(0, clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	messaging.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
}

// clientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func clientStartup2(rsc startupResource, credentials startupCredentials) error {
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

func connectString(url string, credentials startupCredentials) (string, error) {
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
func accessCredentials(msg *messaging.Message) startupCredentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(func() (user string, pswd string, err error)); ok {
			return fn
		}
	}
	return nil
}

// accessResource - access function for a resource in a message
func accessResource(msg *messaging.Message) startupResource {
	if msg == nil || msg.Content == nil {
		return startupResource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(struct{ Uri string }); ok {
			return url
		}
	}
	return startupResource{}
}
