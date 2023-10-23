package cmds

import (
	"context"
	"os"

	"github.com/aserto-dev/azm/migrate"
	client "github.com/aserto-dev/go-aserto/client"
)

type MigrateCmd struct {
}

func (a *MigrateCmd) Run(c *Common) error {
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

	if err := m.Process(clnt.Conn); err != nil {
		return err
	}

	if err := m.Normalize(); err != nil {
		return err
	}

	if err := m.Validate(); err != nil {
		return err
	}

	if err := m.Write(os.Stdout); err != nil {
		return err
	}

	return nil
}
