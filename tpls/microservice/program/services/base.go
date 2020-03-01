package services

/* 客户端和管理端公用逻辑 */

// Base 基础服务对象
type Base struct {
}

var (
	// 保证前后台服务使用一个base
	base *Base
)

// NewBase 创建基础服务对象
func NewBase() Base {
	if base != nil {
		return *base
	}
	base = &Base{}
	return *base
}
