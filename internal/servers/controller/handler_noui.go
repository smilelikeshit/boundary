package controller

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/hashicorp/boundary/internal/observability/event"
)

func devPassthroughHandler(passthroughDir string) http.Handler {
	const op = "controller.devPassthroughHandler"
	ctx := context.TODO()
	// Panic may not be ideal but this is never a production call and it'll
	// panic on startup. We could also just change the function to return
	// an error.
	abs, err := filepath.Abs(passthroughDir)
	if err != nil {
		panic(err)
	}
	event.WriteSysEvent(ctx, op, "serving passthrough files at /", "path", abs)
	fs := http.FileServer(http.Dir(abs))
	prefixHandler := http.StripPrefix("/", fs)

	return prefixHandler
}

var handleUi = func(c *Controller) http.Handler {
	if c.conf.RawConfig.PassthroughDirectory != "" {
		return devPassthroughHandler(c.conf.RawConfig.PassthroughDirectory)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
