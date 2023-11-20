package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/content"
	"github.com/advanced-go/messaging/core"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	dbClient  *pgxpool.Pool
	clientLoc = PkgUri + "/Startup"
)

var clientStartup core.MessageHandler = func(msg core.Message) {
	if isStarted() {
		return
	}
	start := time.Now()
	rsc := content.AccessResource(&msg)
	credentials := content.AccessCredentials(&msg)
	err := ClientStartup(rsc, credentials)
	if err != nil {
		core.ReplyTo(msg, runtime.NewStatusError(0, clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	core.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
}

// ClientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func ClientStartup(rsc content.Resource, credentials content.Credentials) error {
	if isStarted() {
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
		ClientShutdown()
		return errors.New(fmt.Sprintf("unable to acquire connection from pool: %v\n", err1))
	}
	conn.Release()
	setStarted()
	return nil
}

func ClientShutdown() {
	if dbClient != nil {
		resetStarted()
		dbClient.Close()
		dbClient = nil
	}
}

func connectString(url string, credentials content.Credentials) (string, error) {
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
