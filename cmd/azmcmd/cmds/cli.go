package cmds

type CLI struct {
	Common
	Query   QueryCmd   `cmd:"" help:"query model metadata"`
	Migrate MigrateCmd `cmd:"" help:"migrate directory v2 metadata to an annotated v3 manifest"`
	Version VersionCmd `cmd:"" help:"version information"`
}
