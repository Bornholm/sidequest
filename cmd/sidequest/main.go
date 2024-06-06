package main

import (
	"log/slog"
	"os"

	"github.com/bornholm/sidequest/internal/env"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/bornholm/sidequest/migrations"
)

func main() {
	isDev := env.Bool("SIDEQUEST_DEV", false)
	dataDir := env.String("SIDEQUEST_DATA_DIR", "./data")
	publicDir := env.String("SIDEQUEST_PUBLIC_DIR", "./dist")

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDev:     isDev,
		DefaultDataDir: dataDir,
	})

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isDev,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDir), true))
		return nil
	})

	if err := app.Start(); err != nil {
		slog.Error("could not start app", slog.AnyValue(errors.WithStack(err)))
		os.Exit(1)
	}
}
