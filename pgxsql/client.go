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
	err := clientStartup2(msg.Config)
	if err != nil {
		messaging.SendReply(msg, runtime.NewStatusError(0, clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	messaging.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
}

// clientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func clientStartup2(cfg *runtime.StringsMap) error {
	if isReady() {
		return nil
	}
	if cfg == nil {
		return errors.New("error: strings map configuration is nil")
	}
	url, status := cfg.Get(uriConfigKey)
	if !status.OK() {
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

func connectString(url string, cfg *runtime.StringsMap) (string, error) {
	user, status := cfg.Get(userConfigKey)
	pswd, status0 := cfg.Get(pswdConfigKey)
	// Username and password can be in the connect string Url
	if !status.OK() && !status0.OK() {
		return url, nil
	}
	if !status.OK() {
		return "", errors.New("error: user is not configured")
	}
	if !status0.OK() {
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
