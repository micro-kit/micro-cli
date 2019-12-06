package microdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/micro-kit/micro-cli/program/common"
)

// MicroDB 操作存储到本地的json db
type MicroDB struct {
	isInit         bool
	Service        *Service     `json:"service"`         // 服务名
	Messages       []*Message   `json:"messages"`        // message - 结构体
	Enums          []*Enums     `json:"enums"`           // 枚举列表
	PackageName    string       `json:"package_name"`    // 包名
	PackageComment string       `json:"package_comment"` // 包注释
	ProjectInfo    *ProjectInfo `json:"project_info"`    // 项目信息

	DbPath string `json:"-"` // 存储项目相关配置目录
}

// ProjectInfo 项目信息 new 指令存储，和浏览器缓存一致
type ProjectInfo struct {
	Name        string `json:"name"`        // 服务名 和文件夹名一致 基本和MicroDB的Service一致
	SrcPath     string `json:"src_path"`    // 源代码相对于GOPATH的路径
	Description string `json:"description"` // 项目描述
}

// NewMicroDB 创建kit db操作对象
func NewMicroDB(dbPath string) *MicroDB {
	return &MicroDB{
		isInit:      false,
		Service:     new(Service),
		Messages:    make([]*Message, 0),
		Enums:       make([]*Enums, 0),
		ProjectInfo: new(ProjectInfo),
		DbPath:      dbPath,
	}
}

// InitForProto 初始化kit db
func (db *MicroDB) InitForProto(protoFile string) (err error) {
	if protoFile == "" {
		return errors.New("proto文件路径不能为空")
	}
	reader, err := os.Open(protoFile)
	if err != nil {
		// log.Println(err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		// log.Println(err)
		return
	}
	// 已初始化
	db.isInit = true
	// 解析包名
	for _, each := range definition.Elements {
		if s, ok := each.(*proto.Package); ok {
			// log.Println(p.Name)
			db.initPackageName(s)
		}
	}
	if packageName, _ := db.GetPackageName(); packageName == "" {
		return errors.New("包名未定义")
	}

	proto.Walk(definition,
		proto.WithEnum(db.initEnums),
		proto.WithMessage(db.initMessage),
		proto.WithService(db.initService))
	// 处理rpc，不能使用proto.WithRPC(db.initRPC)
	for _, each := range definition.Elements {
		if s, ok := each.(*proto.Service); ok {
			for _, v := range s.Elements {
				if s1, ok := v.(*proto.RPC); ok {
					db.initRPC(s1)
				}
			}
		}
	}

	return nil
}

// GetService 获取Service
func (db *MicroDB) GetService() (*Service, error) {
	if db.isInit == false {
		return nil, errors.New("kit db未初始化")
	}
	return db.Service, nil
}

// GetMessages 获取Service
func (db *MicroDB) GetMessages() ([]*Message, error) {
	if db.isInit == false {
		return nil, errors.New("kit db未初始化")
	}
	return db.Messages, nil
}

// GetEnums 获取enums
func (db *MicroDB) GetEnums() ([]*Enums, error) {
	if db.isInit == false {
		return nil, errors.New("kit db未初始化")
	}
	return db.Enums, nil
}

// 初始化Service
func (db *MicroDB) initService(service *proto.Service) {
	if service == nil {
		log.Println("service is nil")
		return
	}
	// 服务信息
	db.Service = &Service{
		Name:     service.Name,
		Rpcs:     make([]*RPC, 0),
		Comment:  db.GetComment(service.Comment),
		Position: service.Position,
	}
	return
}

// 初始化RPC
func (db *MicroDB) initRPC(rpc *proto.RPC) {
	// log.Println("initRPC")
	if rpc == nil {
		log.Println("rpc is nil")
		return
	}
	request, _ := db.GetMessageByName(rpc.RequestType)
	returns, _ := db.GetMessageByName(rpc.ReturnsType)
	// 服务信息
	db.Service.Rpcs = append(db.Service.Rpcs, &RPC{
		Name:        rpc.Name,
		RequestType: rpc.RequestType,
		ReturnType:  rpc.ReturnsType,
		Request:     request,
		Return:      returns,
		Comment:     db.GetComment(rpc.Comment),
		Position:    rpc.Position,
	})
	return
}

// 初始化Message列表
func (db *MicroDB) initMessage(message *proto.Message) {
	// log.Println("initMessage")
	if message == nil {
		log.Println("message is nil")
		return
	}
	// 包名
	packageName, err := db.GetPackageName()
	if err != nil {
		log.Println(err)
		return
	}
	// 获取字段列表
	fields := make([]*MessageField, 0)
	for _, v := range message.Elements {
		normalField, ok := v.(*proto.NormalField)
		if ok == true {
			fieldType := normalField.Type
			if normalField.Repeated == true {
				if common.IsBasicType(fieldType) == true {
					fieldType = "[]" + fieldType
				} else {
					fieldType = "[]*" + packageName + "." + fieldType
				}
			}
			fields = append(fields, &MessageField{
				Name:    common.StrFirstToUpper(normalField.Name, true),
				Type:    fieldType,
				Number:  normalField.Sequence,
				Comment: db.GetComment(normalField.Comment),
			})
			continue
		}
		mapField, ok := v.(*proto.MapField)
		if ok == true {
			fieldType := ""
			// 判断map key是不是基础数据类型
			if common.IsBasicType(mapField.KeyType) == true {
				fieldType = fmt.Sprintf("map[%s]", mapField.KeyType)
			} else {
				fieldType = fmt.Sprintf("map[*%s]", packageName+"."+mapField.KeyType)
			}
			// 判断值类型是不是基础数据类型
			if common.IsBasicType(mapField.Type) == true {
				fieldType = fieldType + mapField.Type
			} else {
				fieldType = fieldType + "*" + packageName + "." + mapField.Type
			}

			fields = append(fields, &MessageField{
				Name:    common.StrFirstToUpper(mapField.Name, true),
				Type:    fieldType,
				Number:  mapField.Sequence,
				Comment: db.GetComment(mapField.Comment),
			})
		}
		oneOfField, ok := v.(*proto.OneOfField)
		if ok == true {
			fields = append(fields, &MessageField{
				Name:    common.StrFirstToUpper(oneOfField.Name, true),
				Type:    oneOfField.Type,
				Number:  oneOfField.Sequence,
				Comment: db.GetComment(oneOfField.Comment),
			})
		}
	}
	db.Messages = append(db.Messages, &Message{
		Name:    common.StrFirstToUpper(message.Name, true),
		Fields:  fields,
		Comment: db.GetComment(message.Comment),
	})
	// fmt.Println(db.Messages)
	return
}

// GetMessageByName 根据 message name 获取 Message
func (db *MicroDB) GetMessageByName(name string) (msg *Message, err error) {
	if db.isInit == false {
		return nil, errors.New("kit db未初始化")
	}
	// js, _ := json.Marshal(db.Messages)
	// fmt.Println(string(js))
	// fmt.Println(name)
	for _, v := range db.Messages {
		if v.Name == common.StrFirstToUpper(name, true) {
			return v, nil
		}
	}
	return
}

// 初始化Enums列表
func (db *MicroDB) initEnums(enum *proto.Enum) {
	if enum == nil {
		log.Println("enum is nil")
		return
	}
	fields := make([]*EnumsField, 0)

	if len(enum.Elements) > 0 {
		for _, vv := range enum.Elements {
			field, ok := vv.(*proto.EnumField)
			if ok == true {
				fields = append(fields, &EnumsField{
					Name:    common.StrFirstToUpper(field.Name, true),
					Value:   field.Integer,
					Comment: db.GetComment(field.Comment),
				})
			}

		}
	}
	db.Enums = append(db.Enums, &Enums{
		Name:    common.StrFirstToUpper(enum.Name, true),
		Fields:  fields,
		Comment: db.GetComment(enum.Comment),
	})
	return
}

// GetEnumsByName 根据 enums name 获取 Enums
func (db *MicroDB) GetEnumsByName(name string) (msg *Enums, err error) {
	if db.isInit == false {
		return nil, errors.New("kit db未初始化")
	}
	for _, v := range db.Enums {
		if v.Name == common.StrFirstToUpper(name, true) {
			return v, nil
		}
	}
	return
}

// GetServiceJSONString 获取service json 字符串
func (db *MicroDB) GetServiceJSONString() (string, error) {
	if db.isInit == false {
		return "", errors.New("kit db未初始化")
	}
	js, err := json.Marshal(db.Service)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

// 初始化packageName
func (db *MicroDB) initPackageName(pkg *proto.Package) {
	if pkg == nil {
		log.Println("proto 未定义任何 package")
		return
	}
	db.PackageName = pkg.Name
	db.PackageComment = db.GetComment(pkg.Comment)
	return
}

// GetPackageName 获取包名
func (db *MicroDB) GetPackageName() (string, error) {
	if db.isInit == false {
		return "", errors.New("kit db未初始化")
	}
	return db.PackageName, nil
}

// GetComment 获取注释内容
func (db *MicroDB) GetComment(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	if len(comment.Lines) > 0 {
		msgs := make([]string, 0)
		for _, v := range comment.Lines {
			msgs = append(msgs, "//"+v)
		}
		return strings.Join(msgs, "\n")
	} else {
		return comment.Message()
	}
}

// InRpcs 检查rpc方法是否已经存在
func (db *MicroDB) InRpcs(name string) bool {
	if db == nil || db.Service == nil || len(db.Service.Rpcs) == 0 {
		return false
	}
	for _, rpc := range db.Service.Rpcs {
		if name == rpc.Name {
			return true
		}
	}
	return false
}
