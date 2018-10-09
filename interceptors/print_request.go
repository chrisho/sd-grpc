package interceptors

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"encoding/json"
	"google.golang.org/grpc/grpclog"
)

func PrintRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) error {

	contents := make(map[string]string)

	contents["grpc_method"] = info.FullMethod
	contents["grpc_param"] = fmt.Sprint(req)

	contentsJson, err := json.Marshal(contents)
	if err != nil {
		return err
	}
	grpclog.Infoln(contentsJson)
	return nil
}