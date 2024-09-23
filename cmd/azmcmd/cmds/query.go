package cmds

import (
	"fmt"
	"os"

	"github.com/aserto-dev/azm/cmd/azmcmd/pkg/table"
	v3 "github.com/aserto-dev/azm/v3"
)

type QueryCmd struct {
	Query    string `arg:"" `
	Filename string `flag:"" short:"f" name:"file" default:"manifest.yaml"`
}

func (cmd *QueryCmd) Run(c *Common) error {
	mod, err := v3.LoadFile(cmd.Filename)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "model version %d\n", mod.Version)

	tab := table.New(os.Stdout).WithColumns("ObjectType", "Relation", "SubjectType")

	for objTypeName, objType := range mod.Objects {
		for relTypeName, relType := range objType.Relations {
			for _, subTypeName := range relType.SubjectTypes {
				tab.WithRow(string(objTypeName), string(relTypeName), string(subTypeName))
			}
		}

		for relTypeName, relType := range objType.Permissions {
			for _, subTypeName := range relType.SubjectTypes {
				tab.WithRow(string(objTypeName), string(relTypeName), string(subTypeName))
			}
		}

	}

	tab.Do()

	return nil
}
