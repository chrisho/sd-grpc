package sdgrpc

import (
	"github.com/chrisho/sd-grpc/interceptors"
	"github.com/chrisho/sd-helper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"net"
	"strings"
)

type sdgrpc struct {
	Server *grpc.Server
	opts []grpc.ServerOption
}

func NewSdGrpc() (*sdgrpc, error) {
	sg := &sdgrpc{}
    err := sg.configCredentials()
    if err != nil {
		return nil, err
	}

	// 注册interceptor
	sg.opts = append(sg.opts, grpc.UnaryInterceptor(sg.getInterceptor()))

	// 实例化服务
	sg.Server = grpc.NewServer(sg.opts...)

	return sg, nil
}

// 注册协议
func (s *sdgrpc) PbRegister(register func(*grpc.Server, interface{}), srv interface{}) {
	register(s.Server, srv)
}

// 启动服务
func (s *sdgrpc) Run() error {

	debug() // 根据配置设置debug

	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpclog.Infoln("listening TCP " + port)
	err = s.Server.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}


// 配置证书
func (s *sdgrpc) configCredentials() error {
	if strings.ToLower(sdhelper.GetEnv("SSL")) != "true" {
		return nil
	}

	certFile := sdhelper.GetEnv("SSLCertFile")
	keyFile := sdhelper.GetEnv("SSLKeyFile")
	transportCrd, err := credentials.NewServerTLSFromFile(path+"/"+certFile, path+"/"+keyFile)
	if err != nil {
		return err
	}

	s.opts = append(s.opts, grpc.Creds(transportCrd))

	return nil
}

// 获取拦截器方法
func (s *sdgrpc) getInterceptor() func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

			err = interceptors.PrintRequest(ctx, req, info)
			if err != nil {
				return nil, err
			}
			// 继续处理请求
			return handler(ctx, req)
	}
}