// {{ .RpcName }} {{ .Comment }}
func (s *{{ .RpcType }}) {{ .RpcName }}(ctx context.Context, req *{{ .BaseServiceNameNotLine }}pb.{{ .RpcName }}Request) (*{{ .BaseServiceNameNotLine }}pb.{{ .RpcName }}Reply, error) {
    // 验证参数是否错误

    // TODO 逻辑代码

    // 返回结果
	reply := &{{ .BaseServiceNameNotLine }}pb.{{ .RpcName }}Reply{
	}
	return reply, nil
}
