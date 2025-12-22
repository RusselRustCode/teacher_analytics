package grpc

import (
	"fmt"
	"net"
	"google.golang.org/grpc"
	"github.com/RusselRustCode/teacher_analytics/core-service/proto"
)

func StartGRPCServer(port string, handler proto.AnalyticsServiceServer) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterAnalyticsServiceServer(s, handler)

	fmt.Printf("gRPC server listening at %v\n", lis.Addr())
	return s.Serve(lis)
}