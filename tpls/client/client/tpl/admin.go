package {{ .BaseServiceNameNotLine }}

import (
	{{ .BaseServiceNameNotLine }}pb "{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
	"github.com/micro-kit/microkit/client"
	"google.golang.org/grpc"
)

var (
	svcAdminName = "{{ .BaseServiceNameNotLine }}"
)

// NewAdminClient 创建管理端客户端
func NewAdminClient() ({{ .BaseServiceNameNotLine }}AdminClient {{ .BaseServiceNameNotLine }}pb.Admin{{ .BaseServiceNameHump }}Client, err error) {
	c, err := client.NewDefaultClient(client.ServiceName(svcAdminName))
	if err != nil {
		return
	}
	// 连接服务端
	err = c.Dial(func(cc *grpc.ClientConn) {
		{{ .BaseServiceNameNotLine }}AdminClient = {{ .BaseServiceNameNotLine }}pb.New{{ .BaseServiceNameHump }}Client(cc)
	})
	return
}
