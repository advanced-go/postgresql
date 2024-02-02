package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
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
	err := clientStartup2(msg.Config)
	if err != nil {
		messaging.SendReply(msg, messaging.NewStatusDurationError(http.StatusOK, time.Since(start), err))
		return
	}
	messaging.SendReply(msg, messaging.NewStatusDuration(http.StatusOK, time.Since(start)))
}

// clientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func clientStartup2(cfg map[string]string) error {
	if isReady() {
		return nil
	}
	if cfg == nil {
		return errors.New("error: strings map configuration is nil")
	}
	url, ok := cfg[uriConfigKey]
	if !ok {
		return errors.New("database URL is empty")
	}
	// Create connection string with credentials
	s, err := connectString(url, cfg)
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

func connectString(url string, cfg map[string]string) (string, error) {
	user, ok := cfg[userConfigKey]
	pswd, ok0 := cfg[pswdConfigKey]
	// Username and password can be in the connect string Url
	if !ok && !ok0 {
		return url, nil
	}
	if !ok {
		return "", errors.New("error: user is not configured")
	}
	if !ok0 {
		return "", errors.New("error: password is not configured")
	}
	return fmt.Sprintf(url, user, pswd), nil
}

/*
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


*/
