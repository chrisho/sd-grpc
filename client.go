package sdgrpc

import (
	"os"
	"strings"

	"github.com/chrisho/sd-helper/stringx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type client struct {
	conn *grpc.ClientConn
	opts []grpc.DialOption
	optsCallOption []grpc.CallOption
}

func NewClientLocal(serviceName string) (*grpc.ClientConn, error) {

	debug() // 根据配置设置debug

	os.Setenv("ClientConn", "1") // 配置本地ip，测试用
	return NewClient(serviceName)
}

func NewClient(serviceName string) (*grpc.ClientConn, error) {
	serviceName = stringx.ConvertUnderlineToWhippletree(serviceName)
	host := serviceName + os.Getenv("SSLSuffixServerName")

	var address string
	if os.Getenv("ClientConn") == "1" {
		address = "127.0.0.1"
	} else {
		address = serviceName + os.Getenv("ClusterSuffixDomain")
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
	if strings.ToLower(os.Getenv("SSL")) != "true" {
		c.opts = append(c.opts, grpc.WithInsecure())
		return nil
	}

	certFile := os.Getenv("SSLCACertFile")
	creds, err := credentials.NewClientTLSFromFile(path+"/"+certFile, host)
	if err != nil {
		return err
	}
	c.opts = append(c.opts, grpc.WithTransportCredentials(creds))
	return nil
}