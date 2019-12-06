package microdb

import "text/scanner"

// pb文件解析相关结构体

// Service 对应一个GRPC服务
type Service struct {
	Name     string `json:"name"`
	Rpcs     []*RPC `json:"rpcs"`
	Comment  string `json:"comment"`
	Position scanner.Position
}

// GetRPCForName 根据name获取service下的一个rpc
func (s *Service) GetRPCForName(name string) *RPC {
	for _, rpc := range s.Rpcs {
		if rpc.Name == name {
			return rpc
		}
	}
	return nil
}

// RPC rpc服务列表
type RPC struct {
	Name        string   `json:"name"`
	RequestType string   `json:"request_type"`
	ReturnType  string   `json:"return_type"`
	Request     *Message `json:"request"`
	Return      *Message `json:"return"`
	Comment     string   `json:"comment"`
	Position    scanner.Position
}

// Message proto文件中的message列表
type Message struct {
	Name    string          `json:"name"`
	Fields  []*MessageField `json:"fields"`
	Comment string          `json:"comment"`
}

// GetMessageFieldForName 根据字段名获取字段嘻嘻
func (m *Message) GetMessageFieldForName(name string) *MessageField {
	for _, v := range m.Fields {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// MessageField proto文件中的message的字段列表
type MessageField struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Number  int    `json:"number"`
	Comment string `json:"comment"`
}

// Enums 枚举
type Enums struct {
	Name    string        `json:"name"`
	Fields  []*EnumsField `json:"fields"`
	Comment string        `json:"comment"`
}

// EnumsField 美剧值列表
type EnumsField struct {
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Comment string `json:"comment"`
}
