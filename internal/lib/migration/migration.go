package migration

import (
	"context"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas/atlasexec"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

func Run(ctx context.Context, postgresDsn string) {
	workdir, err := atlasexec.NewWorkingDir(atlasexec.WithMigrations(os.DirFS("./migrations")))
	if err != nil {
		log.Fatal(err)
	}
	defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: postgresDsn,
	})
	if err != nil {
		log.Fatal(err)
	}

	logger.InfoInConsole(fmt.Sprintf("applied %d migrations", len(res.Applied)))
}
