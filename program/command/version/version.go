package version

import (
	"fmt"

	"github.com/mitchellh/cli"
)

// New version 指令
func New(ui cli.Ui, version, gitHash string) *cmd {
	return &cmd{
		UI:      ui,
		version: version,
		gitHash: gitHash,
	}
}

type cmd struct {
	UI      cli.Ui
	version string
	gitHash string
}

func (c *cmd) Run(_ []string) int {
	c.UI.Output(fmt.Sprintf("micro-cli %s\ngitHash: %s", c.version, c.gitHash))
	return 0
}

func (c *cmd) Synopsis() string {
	return "打印micro-cli程序版本"
}

func (c *cmd) Help() string {
	return ""
}
