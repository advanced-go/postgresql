package pgxsql

import (
	"context"
	"embed"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/controller"
	"io/fs"
)

const (
	controllersPath      = "resource/controllers.json"
	queryControllerName  = "postgresql-query"
	insertControllerName = "postgresql-insert"
	updateControllerName = "postgresql-update"
	deleteControllerName = "postgresql-delete"
)

var (
	//go:embed resource/*
	f  embed.FS
	cm *controller.Map
)

func init() {
	buf, err := fs.ReadFile(f, controllersPath)
	if err != nil {
		fmt.Printf("controller.init(\"%v\") failure: [%v]\n", PkgPath, err)
		return
	}
	cm, err = controller.NewMap(buf)
	if err != nil {
		fmt.Printf("controller.init(\"%v\") failure: [%v]\n", PkgPath, err)
	}
}

func apply(ctx context.Context, newCtx *context.Context, req *request, statusCode access.StatusCodeFunc) func() {
	var c *controller.Controller
	if cm != nil {
		c, _ = cm.Get(req.controllerName)
	}
	if c == nil {
		c = new(controller.Controller)
		c.Name = "error"
		c.Duration = 0
	}
	return controller.Apply(ctx, newCtx, method(req), req.uri, c.Name, req.header, c.Duration, statusCode)
}
