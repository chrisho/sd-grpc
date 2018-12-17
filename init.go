package sdgrpc

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"os"
	"strings"
)

const envFile = "/config/conf.env"

var (
	path string

	port      = ":50051"
	debugPort = ""
)

func init() {
	// 配置日志logger
	grpclog.SetLoggerV2(
		NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard))

	path, _ = os.Getwd()
	err := godotenv.Load(path + envFile)
	if err != nil {
		grpclog.Fatalln(err)
	}
}

func debug() {
	if strings.ToLower(os.Getenv("Debug")) == "true" {
		if os.Getenv("DebugServerListenPort") != "" {
			debugPort = os.Getenv("DebugServerListenPort")
		}
	}
	if debugPort != "" {
		port = debugPort
	}
}