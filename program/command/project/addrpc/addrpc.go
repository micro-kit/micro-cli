package addrpc

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/codeskyblue/go-sh"
	"github.com/joho/godotenv"
	"github.com/micro-kit/micro-cli/program/command/flags"
	"github.com/micro-kit/micro-cli/program/common"
	"github.com/micro-kit/micro-cli/program/microdb"
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

	microDb *microdb.MicroDB // 项目配置db文件
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.serviceName, "svc", "", "服务名")
	c.flags.StringVar(&c.rpcName, "rpc", "", "rpc服务方法名")
	c.flags.StringVar(&c.rpcType, "type", "foreground", "admin | foreground 前台或后台服务方法 默认前台")
	c.flags.StringVar(&c.comment, "comment", "", "注释")
	c.flags.StringVar(&c.rootPath, "root", "", "服务上层目录相对GOPATH路径\n为空取环境变量ROOT_PATH\n默认github.com/micro-kit，$GOPATH/src/$root")                                    // 项目根路径
	c.flags.StringVar(&c.clientRootPath, "croot", "", "服务客户端库相对GOPATH路径\n为空取环境变量MICROKIT_CLIENT_ROOT\n默认github.com/micro-kit/microkit-client，$GOPATH/src/$croot") // 项目根路径

	c.help = flags.Usage(help, c.flags)
}

// 通过环境变量处理默认值
func (c *cmd) initEnv() {
	projectEnvPath := os.Getenv("GOPATH") + "/src/" + strings.Trim(c.rootPath, "/") + "/" + c.serviceName + "-service" + "/.micro-db/.env"
	if ext, _ := common.PathExists(projectEnvPath); ext == true {
		err := godotenv.Load(projectEnvPath)
		if err != nil {
			log.Println(err)
		}
	}
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

// 初始化已存在项目的MicroDB
func (c *cmd) initMicroDB() error {
	projectDBPath := os.Getenv("GOPATH") + "/src/" + strings.Trim(c.rootPath, "/") + "/" + c.serviceName + "-service" + "/.micro-db/"
	c.microDb = microdb.NewMicroDB(projectDBPath)
	protoFile := c.getProtoFile()
	err := c.microDb.InitForProto(protoFile)
	if err != nil {
		log.Println("解析proto文件错误,文件路径:" + protoFile)
		log.Println(err)
		return err
	}

	return nil
}

// 获取ob文件路径
func (c *cmd) getProtoFile() string {
	protoFile := os.Getenv("GOPATH") + "/src/" + c.clientRootPath + "/proto/" + c.serviceName + "pb"
	if c.rpcType == "Admin" {
		protoFile += "/admin.proto"
	} else {
		protoFile += "/foreground.proto"
	}
	return protoFile
}

// 获取pb目录
func (c *cmd) getProtoPath() string {
	return os.Getenv("GOPATH") + "/src/" + c.clientRootPath + "/proto/" + c.serviceName + "pb"
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	// 参数盘点
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

	// 处理环境变量默认值
	c.initEnv()
	// 解析项目pb文件
	c.initMicroDB()
	defer func() {
		err := c.microDb.SaveToFileNotCheck()
		if err != nil {
			log.Println("写db.json错误", err)
		}
	}()

	// 创建服务方法
	code := c.addSvcRpc()
	if code != 0 {
		return code
	}

	// 创建pb文件-并生成go文件
	code = c.addProtoRpc()
	if code != 0 {
		return code
	}

	return 0
}

// addSvcRpc 创建rpc服务方法
func (c *cmd) addSvcRpc() int {
	// 判断方法是否已经存在
	if c.microDb.InRpcs(common.StrFirstToUpper(c.rpcName, true)) == true {
		log.Println("RPC方法已经存在", common.StrFirstToUpper(c.rpcName, true))
		return 1
	}
	// 项目路径
	projectPath := os.Getenv("GOPATH") + "/src/" + strings.Trim(c.rootPath, "/") + "/" + c.serviceName + "-service"
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
	// 写服务rpc方法
	err := c.TplFileNew("rpc/rpc.tpl", serviceFilePath, c.tplData)
	if err != nil {
		log.Println("解析模版并写入服务文件错误", err)
		return 1
	}
	return 0
}

// TplFileNew 替换文件中的变量，写入到对应目录 - 追加写入 line不为空，插入此行之后
func (c *cmd) TplFileNew(inFileName, outFilePath string, data map[string]interface{}, line ...int) (err error) {
	startLine := 0
	if len(line) > 0 {
		startLine = line[0]
	}
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
	buff := bytes.NewBuffer([]byte("\n"))
	err = t.Execute(buff, data)
	if err != nil {
		return err
	}
	defer buff.Reset()
	if startLine == 0 {
		// 打开文件
		outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		defer outFile.Close()
		_, err = outFile.Write(buff.Bytes())
		if err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadFile(outFilePath)
		if err != nil {
			return err
		}
		lineDatas := strings.Split(string(body), "\n")
		afterLines := []string{strings.TrimLeft(buff.String(), "\n")}
		afterLines = append(afterLines, lineDatas[startLine:]...)
		lineDatas = append(lineDatas[:startLine], afterLines...)
		outBody := strings.Join(lineDatas, "\n")
		err = ioutil.WriteFile(outFilePath, []byte(outBody), 0755)
		if err != nil {
			return err
		}
	}
	return
}

// AddProtoRpc 添加pb文件rpc方法
func (c *cmd) addProtoRpc() int {
	// 追加写入pb message
	err := c.TplFileNew("rpc/pbmessage.tpl", c.getProtoFile(), c.tplData)
	if err != nil {
		log.Println("解析pb message存入文件错误", err)
		return 1
	}
	// 将rpc定义插入services
	startLine := 0
	if c.microDb != nil && c.microDb.Service != nil {
		if len(c.microDb.Service.Rpcs) > 0 {
			startLine = c.microDb.Service.Rpcs[len(c.microDb.Service.Rpcs)-1].Position.Line
		} else {
			startLine = c.microDb.Service.Position.Line
		}
	}
	// log.Println("rpc插入位置", startLine)
	if startLine == 0 {
		log.Println("解析出pb信息中不存在services定义")
		return 1
	}
	err = c.TplFileNew("rpc/pbrpc.tpl", c.getProtoFile(), c.tplData, startLine)
	if err != nil {
		log.Println("解析pb 将rpc定义插入services错误", err)
		return 1
	}
	// pb文件生成go文件
	err = sh.NewSession().SetDir(c.getProtoPath()).Command("./gen.sh").Run()
	if err != nil {
		log.Println("执行gen.sh错误，生成pb对应go文件错误", err)
		return 1
	}
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const synopsis = "在创建一个rpc服务方法时的服务命令，可以生成pb服务和服务中服务方法，不包含具体参数"
const help = `
Usage: micro-cli addrpc [options]
`
