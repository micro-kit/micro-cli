package services

import (
	"context"
	"encoding/json"

	"github.com/micro-kit/micro-common/cache"
	"github.com/micro-kit/micro-common/microerror"
	"{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
	"{{ .RootPath }}/{{ .ServiceName }}/program/models"
)


/* 提供给管理后台使用的rpc */

// Admin 实现grpc管理端rpc接口
type Admin struct {
	Base
}
