package main

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/azm/cmd/azmcmd/cmds"
)

func main() {
	cli := cmds.CLI{
		Common: cmds.Common{},
	}

	ctx := kong.Parse(&cli,
		kong.Name("azm"),
		kong.Description("Authorization Model Tester"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:   true,
			FlagsLast: true,
		}),
	)

	err := ctx.Run(&cli.Common)
	ctx.FatalIfErrorf(err)
}
