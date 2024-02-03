package pgxsql

import (
	"context"
	"embed"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/controller"
	"github.com/advanced-go/core/runtime"
	"io/fs"
	"net/http"
)

const (
	controllersPath      = "resource/controllers.json"
	controllerLookup     = PkgPath + ":lookupController"
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

func statusCode(s **runtime.Status) access.StatusCodeFunc {
	return func() int {
		if s == nil || *s == nil {
			return http.StatusOK
		}
		return (*(s)).Code
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
