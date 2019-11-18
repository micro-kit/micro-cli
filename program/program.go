package program

import (
	"log"
	"os"

	"github.com/micro-kit/micro-cli/program/command"
	"github.com/mitchellh/cli"
)

/* 程序实例 */

// Program 应用
type Program struct {
	version string // 版本
	gitHash string // git提交id
}

// New 创建实例
func New(version, gitHash string) *Program {
	if version == "" {
		version = "0.0.1"
	}
	p := &Program{
		version: version,
		gitHash: gitHash,
	}

	// 初始化所有指令
	p.RegisterAll()

	return p
}

// Run 运行程序
func (p *Program) Run() int {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "--version" {
			args = []string{"version"}
			break
		}
	}

	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}
	cmds := command.Map(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}
	cli := &cli.CLI{
		Args:         args,
		Commands:     cmds,
		Autocomplete: true,
		Name:         "micro-cli",
		HelpFunc:     cli.FilteredHelpFunc(names, cli.BasicHelpFunc("micro-cli")),
	}

	exitStatus, err := cli.Run()
	if err != nil {
		log.Println(err)
	}
	return exitStatus
}
