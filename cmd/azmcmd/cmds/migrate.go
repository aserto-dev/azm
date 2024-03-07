package cmds

import (
	"context"
	"os"

	"github.com/aserto-dev/azm/migrate"
	client "github.com/aserto-dev/go-aserto/client"
)

type MigrateCmd struct {
	Filename      string `flag:"" name:"file" default:"manifest.yaml"`
	Description   string `flag:"" name:"desc" default:"automatic migration of v2 model to annotated v3 manifest"`
	InclHeader    bool   `flag:"" name:"hdr" default:"false" help:"incl header"`
	InclTimestamp bool   `flag:"" name:"ts" default:"false" help:"incl timestamp"`
}

func (cmd *MigrateCmd) Run(c *Common) error {
	ctx := context.Background()

	opts := []client.ConnectionOption{
		client.WithAddr(c.Host),
		client.WithAPIKeyAuth(c.APIKey),
		client.WithTenantID(c.TenantID),
		client.WithInsecure(c.Insecure),
	}

	clnt, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return err
	}

	m := migrate.NewMigrator()

	if err := m.Load(clnt); err != nil {
		return err
	}

	if err := m.Process(); err != nil {
		return err
	}

	writerOpts := []migrate.WriterOption{
		migrate.WithFilename(cmd.Filename),
		migrate.WithDescription(cmd.Description),
		migrate.WithHeader(cmd.InclHeader),
		migrate.WithTimestamp(cmd.InclTimestamp),
	}

	if err := m.Write(os.Stdout, writerOpts...); err != nil {
		return err
	}

	return nil
}
