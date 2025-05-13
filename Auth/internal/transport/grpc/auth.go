package transport

import (
	"auth/internal/service"
	"auth/pkg/proto"

	"google.golang.org/grpc"
)

func Register(gRPCServer *grpc.Server, auth service.AuthService) {
	proto.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

type serverAPI struct {
	proto.UnimplementedAuthServer
	auth service.AuthService
}
