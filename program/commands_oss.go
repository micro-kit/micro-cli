package program

import (
	"github.com/micro-kit/micro-cli/program/command"
	"github.com/micro-kit/micro-cli/program/command/helm"
	"github.com/micro-kit/micro-cli/program/command/project"
	"github.com/micro-kit/micro-cli/program/command/version"
	"github.com/mitchellh/cli"
)

// RegisterAll 注册命令
func (p *Program) RegisterAll() {
	command.Register("version", func(ui cli.Ui) (cli.Command, error) { return version.New(ui, p.version, p.gitHash), nil })
	command.Register("project", func(ui cli.Ui) (cli.Command, error) { return project.New(ui), nil })
	command.Register("helm", func(ui cli.Ui) (cli.Command, error) { return helm.New(ui), nil })
}
