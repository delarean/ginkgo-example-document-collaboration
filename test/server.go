package collaboration_tests

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"collaboration"
	pb "collaboration/proto"
)

const bufSize = 1024 * 1024

func runServer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	pb.RegisterCollaborationServiceServer(s, collaboration.NewMockCollaborationServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}
