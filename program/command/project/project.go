package project

import (
	"flag"

	"github.com/micro-kit/micro-cli/program/command/flags"
	"github.com/mitchellh/cli"
)

/* 用于项目创建 */

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
	rootPath    string // 项目根路径 - 服务上层文件相对GOPATH路径
	serviceName string // 服务名
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.serviceName, "name", "", "当前服务名")                           // 服务名称
	c.flags.StringVar(&c.rootPath, "root", "", "服务上层文件相对GOPATH路径，$GOPATH/src/$root") // 项目根路径

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

const synopsis = "用于创建项目，客户端代码更新等"
const help = `
Usage: micro-cli project [options]
`
