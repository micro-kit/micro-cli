package addrpc

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/micro-kit/micro-cli/program/command/flags"
	"github.com/micro-kit/micro-cli/program/common"
	"github.com/micro-kit/micro-cli/tpls"
	"github.com/mitchellh/cli"
)

/* 添加rpc服务 */

// New version 指令
func New(ui cli.Ui) *cmd {
	c := &cmd{
		UI: ui,
	}
	c.init()
	return c
}

type cmd struct {
	UI             cli.Ui
	flags          *flag.FlagSet
	help           string                 // 帮助
	serviceName    string                 // 服务名
	rpcName        string                 // rpc服务名
	rpcType        string                 // admin | foreground 前台或后台服务方法 默认前台
	comment        string                 // 函数注释
	rootPath       string                 // 项目根路径 - 服务上层文件相对GOPATH路径
	clientRootPath string                 // 客户端库目录
	tplData        map[string]interface{} // 模版替换参数
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.serviceName, "svc", "", "服务名")
	c.flags.StringVar(&c.rpcName, "rpc", "", "rpc服务方法名")
	c.flags.StringVar(&c.rpcType, "type", "foreground", "admin | foreground 前台或后台服务方法 默认前台")
	c.flags.StringVar(&c.comment, "comment", "", "注释")
	c.flags.StringVar(&c.rootPath, "root", "", "服务上层目录相对GOPATH路径\n为空取环境变量ROOT_PATH\n默认github.com/micro-kit，$GOPATH/src/$root")                                    // 项目根路径
	c.flags.StringVar(&c.clientRootPath, "croot", "", "服务客户端库相对GOPATH路径\n为空取环境变量MICROKIT_CLIENT_ROOT\n默认github.com/micro-kit/microkit-client，$GOPATH/src/$croot") // 项目根路径
	// 处理默认值
	if c.rootPath == "" {
		c.rootPath = os.Getenv("ROOT_PATH")
		if c.rootPath == "" {
			c.rootPath = "github.com/micro-kit"
		}
	}
	if c.clientRootPath == "" {
		c.clientRootPath = os.Getenv("MICROKIT_CLIENT_ROOT")
		if c.clientRootPath == "" {
			c.clientRootPath = "github.com/micro-kit/microkit-client"
		}
	}
	// 去除右侧/
	c.rootPath = strings.TrimRight(c.rootPath, "/")
	c.clientRootPath = strings.TrimRight(c.clientRootPath, "/")

	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	if c.serviceName == "" {
		log.Println("服务名不能为空")
		c.UI.Output(c.Help())
		return 0
	}
	if c.rpcName == "" {
		log.Println("rpc服务方法名不能为空")
		c.UI.Output(c.Help())
		return 0
	}
	if c.rpcType != "foreground" && c.rpcType != "admin" {
		log.Println("rpc方法类型必须是 foreground | admin")
		c.UI.Output(c.Help())
		return 0
	}
	c.rpcName = common.StrFirstToUpper(c.rpcName, true)
	c.rpcType = common.StrFirstToUpper(c.rpcType, true)
	// 创建服务方法
	code := c.addSvcRpc()
	if code != 0 {
		return code
	}

	return 0
}

// addSvcRpc 创建rpc服务方法
func (c *cmd) addSvcRpc() int {
	// 项目路径
	projectPath := os.Getenv("GOPATH") + "/src/" + strings.Trim(c.rootPath, "/") + "/" + c.serviceName
	log.Println("项目路径", projectPath)
	servicePath := projectPath + "/program/services"
	var serviceFilePath string
	if c.rpcType == "Admin" {
		serviceFilePath = servicePath + "/admin.go"
	} else {
		serviceFilePath = servicePath + "/foreground.go"
	}
	// 模版参数
	c.tplData = map[string]interface{}{
		"BaseServiceNameNotLine": strings.ReplaceAll(c.serviceName, "-", ""),
		"RpcName":                common.StrFirstToUpper(c.rpcName, true),
		"Comment":                c.comment,
		"RpcType":                common.StrFirstToUpper(c.rpcType, true),
	}
	err := c.TplFileNew("rpc/rpc.tpl", serviceFilePath, c.tplData)
	if err != nil {
		log.Println("解析模版并写入服务文件错误", err)
		return 1
	}
	return 0
}

// TplFileNew 替换文件中的变量，写入到对应目录
func (c *cmd) TplFileNew(inFileName, outFilePath string, data map[string]interface{}) (err error) {
	// 获取模版文件名
	inFileInfo, err := tpls.AssetInfo(inFileName)
	if err != nil {
		return
	}
	inFName := inFileInfo.Name()
	// 读取模版内容
	tplBytes, err := tpls.Asset(inFileName)
	if err != nil {
		return
	}
	// outFileName := outFilePath + "/" + inFName
	loggingTpl := template.New(inFName)
	t, err := loggingTpl.Parse(string(tplBytes))
	if err != nil {
		return
	}
	// log.Println(outFilePath)
	// 打开文件
	outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer outFile.Close()
	err = t.Execute(outFile, data)
	if err != nil {
		return err
	}
	return
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const synopsis = "在创建一个rpc服务方法时的服务命令，可以生成pb服务和服务中服务方法，不包含具体参数"
const help = `
Usage: micro-cli project addsvc [options]
`
