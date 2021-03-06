package sdgrpc

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func Error(c codes.Code, msg string) error {
	_, file, line, _ := runtime.Caller(1)

	errorLog := fmt.Sprintf("error file : %v ( code : %d) ( line : %d ) \n error info : %v", file, int(c), line, msg)
	fmt.Println(errorLog)

	// 错误代码 200 、 422 、自定义验证错误（10000+）， 不记录日志
	if c != 200 && c != 422 && c < 9000 {
		logrus.Error(errorLog)
	}

	return status.Errorf(c, msg)
}

func ErrorCode(err error) codes.Code {
	s, ok := status.FromError(err)

	if ok {
		return s.Code()
	}
	return 0
}
