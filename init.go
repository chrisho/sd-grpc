package sdgrpc

import (
	"os"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/grpclog"
	"strings"
	"github.com/chrisho/sd-helper"
	"io/ioutil"
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
	if strings.ToLower(sdhelper.GetEnv("Debug")) == "true" {
		if sdhelper.GetEnv("DebugServerListenPort") != "" {
			debugPort = sdhelper.GetEnv("DebugServerListenPort")
		}
	}
	if debugPort != "" {
		port = debugPort
	}
}