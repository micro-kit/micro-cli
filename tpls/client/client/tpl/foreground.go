package {{ .BaseServiceNameNotLine }}

import (
	{{ .BaseServiceNameNotLine }}pb "{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
	"github.com/micro-kit/microkit/client"
	"google.golang.org/grpc"
)

var (
	svcName = "{{ .BaseServiceNameNotLine }}"
)

// NewClient 创建客户端
func NewClient() ({{ .BaseServiceNameNotLine }}Client {{ .BaseServiceNameNotLine }}pb.{{ .BaseServiceNameHump }}Client, err error) {
	c, err := client.NewDefaultClient(client.ServiceName(svcName))
	if err != nil {
		return
	}
	// 连接服务端
	err = c.Dial(func(cc *grpc.ClientConn) {
		{{ .BaseServiceNameNotLine }}Client = {{ .BaseServiceNameNotLine }}pb.New{{ .BaseServiceNameHump }}Client(cc)
	})
	return
}
