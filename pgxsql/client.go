package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sre/core/runtime"
	"github.com/go-sre/host/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var dbClient *pgxpool.Pool
var clientLoc = pkgPath + "/client"

var clientStartup messaging.MessageHandler = func(msg messaging.Message) {
	if IsStarted() {
		return
	}
	start := time.Now()
	db := messaging.AccessDatabaseUrl(&msg)
	credentials := messaging.AccessCredentials(&msg)
	err := ClientStartup(db, credentials)
	if err != nil {
		messaging.ReplyTo(msg, runtime.NewStatusError(clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	messaging.ReplyTo(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
}

// ClientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func ClientStartup(db messaging.DatabaseUrl, credentials messaging.Credentials) error {
	if IsStarted() {
		return nil
	}
	if db.Url == "" {
		return errors.New("database URL is empty")
	}
	// Create connection string with credentials
	s, err := connectString(db.Url, credentials)
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

func connectString(url string, credentials messaging.Credentials) (string, error) {
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
