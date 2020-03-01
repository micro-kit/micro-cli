package project

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

/* 用于项目创建 */

var (
	microServiceTplPath = "microservice" // 默认模版路径
	microClientTplPath  = "client"       // 客户端路径
	clientFiles         = map[string][]string{
		"root": []string{"README.md"},
		"client": []string{
			"admin.go",
			"foreground.go",
			"README.md",
		},
		"proto": []string{
			"admin.proto",
			"foreground.proto",
			"gen.sh",
			"README.md",
		},
	}
	excludedPathNames = map[string]bool{ // 排除的文件列表
		"Desktop.ini": true,
		".DS_Store":   true,
		".git":        true,
	}
)

// New version 指令
func New(ui cli.Ui) *cmd {
	c := &cmd{
		UI: ui,
	}
	c.init()
	return c
}

type cmd struct {
	UI              cli.Ui
	flags           *flag.FlagSet
	help            string                 // 帮助
	rootPath        string                 // 项目根路径 - 服务上层文件相对GOPATH路径
	clientRootPath  string                 // 客户端库目录
	serviceName     string                 // 服务名-此名是输入参数拼接-service
	serviceDesc     string                 // 服务文本描述，介绍功能
	baseServiceName string                 // 输入参数服务名
	tplData         map[string]interface{} // 模版替换参数
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.serviceName, "name", "", "服务名")                                                                                                          // 服务名称
	c.flags.StringVar(&c.serviceDesc, "desc", "", "服务描述")                                                                                                         // 服务描述
	c.flags.StringVar(&c.rootPath, "root", "", "服务上层目录相对GOPATH路径\n为空取环境变量ROOT_PATH\n默认github.com/micro-kit，$GOPATH/src/$root")                                    // 项目根路径
	c.flags.StringVar(&c.clientRootPath, "croot", "", "服务客户端库相对GOPATH路径\n为空取环境变量MICROKIT_CLIENT_ROOT\n默认github.com/micro-kit/microkit-client，$GOPATH/src/$croot") // 项目根路径

	c.help = flags.Usage(help, c.flags)
}

// 通过环境变量处理默认值
func (c *cmd) initEnv() {
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
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	// 处理环境变量默认值
	c.initEnv()

	// 服务名不能为空
	if c.serviceName == "" {
		log.Println("服务名不能为空")
		c.UI.Output(c.Help())
		return 0
	}
	c.baseServiceName = c.serviceName // 保存原服务名
	c.serviceName = c.serviceName + "-service"
	// 创建项目
	code := c.createProject()
	if code != 0 {
		return code
	}

	return 0
}

// 创建项目
func (c *cmd) createProject() int {
	// 项目路径
	projectPath := os.Getenv("GOPATH") + "/src/" + strings.Trim(c.rootPath, "/") + "/" + c.serviceName
	log.Println("项目路径", projectPath)
	if isExis, _ := common.PathExists(projectPath); isExis == false {
		err := os.MkdirAll(projectPath, os.ModePerm)
		if err != nil {
			log.Println("创建服务目录错误", err)
			return 1
		}
	}
	// 检查目录是否是空
	lss, err := common.LsPath(projectPath)
	if err != nil {
		log.Println("查看目录是否为空错误", err)
		return 1
	}
	if len(lss) > 0 {
		log.Println("目录不为空，请在空目录执行project指令，或执行其他指令进行其他操作")
		return 1
	}
	// 模版参数
	c.tplData = map[string]interface{}{
		"BaseServiceNameNotLine": strings.ReplaceAll(c.baseServiceName, "-", ""),
		"BaseServiceNameHump":    common.StrFirstToUpper(c.baseServiceName, true),
		"ServiceName":            c.serviceName,
		"ServiceNameHump":        common.StrFirstToUpper(c.serviceName, true),
		"RootPath":               c.rootPath,
		"MicroKitClientRoot":     c.clientRootPath,
	}
	// 生成文件 - 微服务
	err = c.TreeMicroFilePath(microServiceTplPath)
	if err != nil {
		log.Println("写服务文件错误", err)
		return 1
	}

	// 生成文件 - 客户端
	err = c.TreeClientFilePath(microClientTplPath)
	if err != nil {
		log.Println("写客户端件错误", err)
		return 1
	}

	return 0
}

// LsPath 当前目录下的文件列表
func (c *cmd) LsPath(filePath string) (pathList []string, err error) {
	pathList, err = tpls.AssetDir(filePath)
	if err == nil {
		for k, v := range pathList {
			if _, ok := excludedPathNames[v]; ok == true {
				pathList = append(pathList[:k], pathList[k+1:]...)
			}
		}
	}
	return
}

// TreeMicroFilePath 递归模版目录
func (c *cmd) TreeMicroFilePath(rootPath string) (err error) {
	if rootPath == "" {
		rootPath = microServiceTplPath
	}
	// 获取文件列表
	pathList, err := c.LsPath(rootPath)
	// log.Println(pathList)
	if err != nil {
		if strings.Contains(err.Error(), "not found") == true {
			err = nil
		}
		return
	}
	// 判断路径是目录还是文件-是文件则进行替换-是目录则递归查看
	for _, onePath := range pathList {
		// 完整路径
		completePath := rootPath + "/" + onePath
		tplSubPaths := strings.Split(completePath, "/")
		completeOutPath := os.Getenv("GOPATH") + "/src/" + c.rootPath + "/" + c.serviceName + "/" + strings.Join(tplSubPaths[1:], "/")
		log.Println("创建文件或目录", completeOutPath)
		// continue
		// 判断是路径还是文件
		_, _err := tpls.AssetDir(completePath)
		// _err == nil 表示是目录
		if _err == nil {
			// 判断路径是否存在
			isExist, err := common.PathExists(completeOutPath)
			if err != nil {
				return err
			}
			if isExist == false {
				err = os.MkdirAll(completeOutPath, os.ModePerm)
				if err != nil {
					log.Println("创建目录错误", err)
					return err
				}
			}
			err = c.TreeMicroFilePath(completePath)
			if err != nil {
				return err
			}
		} else {
			err = c.TplFileNew(completePath, completeOutPath, c.tplData)
			if err != nil {
				log.Println("创建文件错误", err)
				return
			}
		}
	}
	return
}

// TreeClientFilePath 递归模版目录
func (c *cmd) TreeClientFilePath(rootPath string) (err error) {
	if rootPath == "" {
		rootPath = microClientTplPath
	}
	// 获取文件列表
	pathList, err := c.LsPath(rootPath)
	// log.Println(pathList)
	if err != nil {
		if strings.Contains(err.Error(), "not found") == true {
			err = nil
		}
		return
	}
	if len(pathList) == 0 {
		log.Println("客户端模版不存在", rootPath)
		return nil
	}
	// 客户端根路径
	clientRootPath := os.Getenv("GOPATH") + "/src/" + c.clientRootPath + "/"
	tplRootPath := rootPath + "/"
	// 创建服务目录
	for k, v := range clientFiles {
		if k == "root" {
			for _, name := range v {
				clientRootPathMd := clientRootPath + name
				isExists, _ := common.PathExists(clientRootPathMd)
				if isExists == false {
					c.TplFileNew(tplRootPath+name, clientRootPathMd, c.tplData)
				}
			}
		} else {
			inFileNameParent := tplRootPath + k + "/tpl/"
			clientRootPathMdParent := clientRootPath + k + "/" + strings.ReplaceAll(c.baseServiceName, "-", "")
			log.Println(inFileNameParent)
			if k == "proto" {
				clientRootPathMdParent = clientRootPath + k + "/" + strings.ReplaceAll(c.baseServiceName, "-", "") + "pb"
			}

			// 目录不存在则创建目录
			isExists, _ := common.PathExists(clientRootPathMdParent)
			if isExists == false {
				err = os.MkdirAll(clientRootPathMdParent, os.ModePerm)
				if err != nil {
					log.Println("创建目录错误", clientRootPathMdParent, err)
					return err
				}
			}
			clientRootPathMdParent += "/"
			// 写文件
			for _, name := range v {
				inFileName := inFileNameParent + name
				clientRootPathMd := clientRootPathMdParent + name
				c.TplFileNew(inFileName, clientRootPathMd, c.tplData)
			}

		}
	}

	return
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

const synopsis = "用于创建项目，客户端代码更新等"
const help = `
Usage: micro-cli project -name [options]
`
