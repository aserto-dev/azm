package cmds

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	ver    string //nolint:gochecknoglobals // set by linker
	date   string //nolint:gochecknoglobals // set by linker
	commit string //nolint:gochecknoglobals // set by linker
)

// Info - version info.
type Info struct {
	Version string
	Date    string
	Commit  string
}

type VersionCmd struct {
}

func (a *VersionCmd) Run(c *Common) error {
	fmt.Fprint(os.Stdout, GetInfo().String())
	return nil
}

// GetInfo - get version stamp information.
func GetInfo() Info {
	if ver == "" {
		ver = "0.0.0"
	}

	if date == "" {
		date = time.Now().Format(time.RFC3339)
	}

	if commit == "" {
		commit = "undefined"
	}

	return Info{
		Version: ver,
		Date:    date,
		Commit:  commit,
	}
}

// String() -- return version info string.
func (vi Info) String() string {
	return fmt.Sprintf("%s g%s %s-%s [%s]",
		vi.Version,
		vi.Commit,
		runtime.GOOS,
		runtime.GOARCH,
		vi.Date,
	)
}
