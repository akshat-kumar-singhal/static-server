package main

import (
	"net/http"
	"strconv"

	"gofr.dev/pkg/gofr"

	"zop.dev/static-server/internal/config"
)

const defaultStaticFilePath = `./static`
const defaultExtensionHTML = ".html"
const indexHTML = "/index.html"
const rootPath = "/"

func main() {
	app := gofr.New()

	staticFilePath := app.Config.GetOrDefault("STATIC_DIR_PATH", defaultStaticFilePath)
	spaMode, _ := strconv.ParseBool(app.Config.GetOrDefault("SPA_MODE", "false"))

	defaultExtension := app.Config.GetOrDefault("DEFAULT_EXTENSION", defaultExtensionHTML)

	handler := &staticFileHandler{
		staticFilePath:   staticFilePath,
		spaMode:          spaMode,
		defaultExtension: defaultExtension,
	}

	app.OnStart(func(ctx *gofr.Context) error {
		handler.fs = ctx.File

		if err := config.HydrateFile(ctx.File, app.Config); err != nil {
			ctx.Logger.Error(err.Error())
		}

		return nil
	})

	app.UseMiddleware(func(next http.Handler) http.Handler {
		handler.next = next
		return handler
	})

	app.AddStaticFiles("/", staticFilePath)

	app.Run()
}
