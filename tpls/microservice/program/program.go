package program

import (
	"log"

	"github.com/micro-kit/micro-common/logger"
	"{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
	"{{ .RootPath }}/{{ .ServiceName }}/program/services"
	"github.com/micro-kit/microkit/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Program 应用实体
type Program struct {
	srv    *server.Server
	logger *zap.SugaredLogger
}

// New 创建应用
func New() *Program {
	// 使用默认服务，如果自定义可设置对应参数
	srv, err := server.NewDefaultServer()
	if err != nil {
		log.Fatalln("创建grpc服务错误", err)
	}
	return &Program{
		srv:    srv,
		logger: logger.Logger,
	}
}

// Run 运行程序
func (p *Program) Run() {
	p.srv.Serve(func(grpcServer *grpc.Server) {
		{{ .BaseServiceNameNotLine }}pb.Register{{ .BaseServiceNameHump }}Server(grpcServer, new(services.Foreground))
	})
	return
}

// Stop 程序结束要做的事
func (p *Program) Stop() {
	p.srv.Stop()
	return
}
