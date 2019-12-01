package services

import (
	"context"
	"encoding/json"

	"github.com/micro-kit/micro-common/cache"
	"github.com/micro-kit/micro-common/microerror"
	"{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
	"{{ .RootPath }}/{{ .ServiceName }}/program/models"
)

/* 提供给客户端使用的rpc */

// Foreground 实现grpc客户端rpc接口
type Foreground struct {
	Base
}
