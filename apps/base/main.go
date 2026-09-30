package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	dls "dls/dls-core"
	"dls/dls-core/adapters/httpapi"
)

func main() {
	app := pocketbase.New()

	// Enable auto migration when running with "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})
	app.RootCmd.AddCommand(newSeedCommand(app))
	app.RootCmd.AddCommand(newIngestCommand(app))

	// The schema is code, installed on every boot. It is idempotent, so a
	// database that already has it is left alone.
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return dls.Migrate(e.App)
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		httpapi.New(dls.New(e.App).Stats).Mount(e)

		// The prerendered site, when it has been built. Registered last so
		// it only catches what the API did not.
		if dir := publicDir(); dir != "" {
			e.Router.GET("/{path...}", apis.Static(os.DirFS(dir), true))
		}
		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// publicDir is where the built site lives. Absent in development, where Vite
// serves it instead.
func publicDir() string {
	dir := os.Getenv("PB_PUBLIC_DIR")
	if dir == "" {
		dir = "pb_public"
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	if _, err := fs.Stat(os.DirFS(abs), "."); err != nil {
		return ""
	}
	return abs
}
