package helm

import (
	"flag"

	"github.com/micro-kit/micro-cli/program/command/flags"
	"github.com/mitchellh/cli"
)

/* 用于生成项目helm部分 */

// New version 指令
func New(ui cli.Ui) *cmd {
	c := &cmd{
		UI: ui,
	}
	c.init()
	return c
}

type cmd struct {
	UI          cli.Ui
	flags       *flag.FlagSet
	help        string // 帮助
	serviceName string // 服务名
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.serviceName, "name", "", "当前服务名")

	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	if c.serviceName == "" {
		c.UI.Output("服务名不能为空!\n")
		c.UI.Output(c.Help())
		return 0
	}
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const synopsis = "用于生成某个服务的helm配置，可以维护某个项目helm部分代码"
const help = `
Usage: micro-cli helm [options]
`
