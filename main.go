package main

import (
	"context"
	"net"

	auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

//TODO: Read from cobra flags
const grpcServerPort = ":8080"
const xApiKey = "x-api-key"
const authorization = "authorization"
const apiKey = "db3f4d66-9d89-4dd4-a865-d3889f558f4f"
const bearerToken = "Bearer YWxpY2VAZXhhbXBsZS5jb206cGFzc3dvcmQ="

type ExtAuthzServer struct {
	logger *zap.Logger
}

func NewExtAuthzServer() *ExtAuthzServer {
	logger, _ := zap.NewProduction()
	return &ExtAuthzServer{
		logger: logger,
	}
}

func (s *ExtAuthzServer) Check(ctx context.Context, req *auth_v3.CheckRequest) (*auth_v3.CheckResponse, error) {
	s.logger.Info("Got request")
	headers := req.Attributes.Request.Http.GetHeaders()
	s.logger.Info("Incoming request headers", zap.Any("headers", headers))

	// Check if X-API-Key is set
	if val, present := headers[xApiKey]; present {
		if val == apiKey {
			return &auth_v3.CheckResponse{
				Status: &status.Status{
					Code: int32(rpc.OK),
				},
				HttpResponse: &auth_v3.CheckResponse_OkResponse{
					OkResponse: &auth_v3.OkHttpResponse{
						HeadersToRemove: []string{xApiKey},
					},
				},
			}, nil
		}
		return &auth_v3.CheckResponse{
			Status: &status.Status{
				Code: int32(rpc.PERMISSION_DENIED),
			},
			HttpResponse: &auth_v3.CheckResponse_DeniedResponse{
				DeniedResponse: &auth_v3.DeniedHttpResponse{
					Status: &envoy_type_v3.HttpStatus{
						Code: envoy_type_v3.StatusCode_Unauthorized,
					},
					Body: rpc.PERMISSION_DENIED.String(),
				},
			},
		}, nil
	}
	// Check if Authorization header is set
	if val, present := headers[authorization]; present {
		if val == bearerToken {
			return &auth_v3.CheckResponse{
				Status: &status.Status{
					Code: int32(rpc.OK),
				},
				HttpResponse: &auth_v3.CheckResponse_OkResponse{
					OkResponse: &auth_v3.OkHttpResponse{
						HeadersToRemove: []string{authorization},
					},
				},
			}, nil
		}
		return &auth_v3.CheckResponse{
			Status: &status.Status{
				Code: int32(rpc.PERMISSION_DENIED),
			},
			HttpResponse: &auth_v3.CheckResponse_DeniedResponse{
				DeniedResponse: &auth_v3.DeniedHttpResponse{
					Status: &envoy_type_v3.HttpStatus{
						Code: envoy_type_v3.StatusCode_Unauthorized,
					},
					Body: rpc.PERMISSION_DENIED.String(),
				},
			},
		}, nil
	}
	return &auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.PERMISSION_DENIED),
		},
		HttpResponse: &auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_Unauthorized,
				},
				Body: rpc.PERMISSION_DENIED.String(),
			},
		},
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", grpcServerPort)
	if err != nil {
		//TODO: ugly handle it better
		panic(err)
	}
	var grpcServerOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(grpcServerOptions...)
	auth_v3.RegisterAuthorizationServer(grpcServer, NewExtAuthzServer())

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}
