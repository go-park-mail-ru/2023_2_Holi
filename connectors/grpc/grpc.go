package grpc

import (
	logs "2023_2_Holi/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))

func Connect(addr string) *grpc.ClientConn {
	if addr == "" {
		logs.Logger.Info("can`t connect by grpc: address is empty")
		return nil
	}

	grpcConn, err := grpc.Dial(
		addr,
		//grpc.WithInsecure(), // TODO add any helpful options?
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logs.LogFatal(logs.Logger, "grpc", "Connect", err, err.Error())
	}

	// TODO healthcheck?

	logs.Logger.Debug("grpc client :", grpcConn)

	return grpcConn
}
