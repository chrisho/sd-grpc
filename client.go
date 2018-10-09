package sdgrpc

import (
	"strings"

	"github.com/chrisho/sd-helper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type client struct {
	conn *grpc.ClientConn
	opts []grpc.DialOption
	optsCallOption []grpc.CallOption
}

func NewClientLocal(serviceName string) (*grpc.ClientConn, error) {

	debug() // 根据配置设置debug

	sdhelper.SetEnv("ClientConn", "1") // 配置本地ip，测试用
	return NewClient(serviceName)
}

func NewClient(serviceName string) (*grpc.ClientConn, error) {
	serviceName = sdhelper.ConvertUnderlineToWhippletree(serviceName)
	host := serviceName + sdhelper.GetEnv("SSLSuffixServerName")

	var address string
	if sdhelper.GetEnv("ClientConn") == "1" {
		address = "127.0.0.1"
	} else {
		address = serviceName + sdhelper.GetEnv("ClusterSuffixDomain")
	}

	c := &client{}

	// 设置接收最大条数
	c.optsCallOption = append(c.optsCallOption, grpc.MaxCallRecvMsgSize(100*1024*1024))
	c.opts = append(c.opts, grpc.WithDefaultCallOptions(c.optsCallOption...))

	// 配置CA证书
	err := c.configCACertFile(host)
	if err != nil {
		return nil, err
	}

	grpclog.Infoln("Certificate Host: ", host)
	grpclog.Infoln("Connect Server: ", address+port)
	c.conn, err = grpc.Dial(address+port, c.opts...)
	if err != nil {
		return nil, err
	}
	return c.conn, nil
}

func (c *client) Close() {
	c.conn.Close()
}

// 配置CA证书
func (c *client) configCACertFile(host string) error {
	if strings.ToLower(sdhelper.GetEnv("SSL")) != "true" {
		c.opts = append(c.opts, grpc.WithInsecure())
		return nil
	}

	certFile := sdhelper.GetEnv("SSLCACertFile")
	creds, err := credentials.NewClientTLSFromFile(path+"/"+certFile, host)
	if err != nil {
		return err
	}
	c.opts = append(c.opts, grpc.WithTransportCredentials(creds))
	return nil
}