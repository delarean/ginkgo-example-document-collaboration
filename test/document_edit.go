package collaboration_tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "collaboration/proto"
)

var _ = Describe("User edits a document", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		conn   *grpc.ClientConn
		client pb.CollaborationServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		var err error
		conn, err = grpc.DialContext(
			ctx,
			"bufnet",
			grpc.WithContextDialer(runServer()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).NotTo(HaveOccurred())
		client = pb.NewCollaborationServiceClient(conn)
	})

	AfterEach(func() {
		cancel()
		conn.Close()
	})

	When("User sends edits to the document", func() {
		It("must be processed successfully", func() {
			req := &pb.EditRequest{
				DocumentId: "doc1",
				User: &pb.User{
					UserId:      "user1",
					DisplayName: "User One",
				},
				Edit: &pb.Edit{
					UserId:    "user1",
					Position:  0,
					Text:      "Hello, Ginkgo!",
					Type:      pb.EditType_INSERT,
					Timestamp: time.Now().Unix(),
				},
			}
			resp, err := client.SendEdit(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Success).To(BeTrue())
		})
	})
})
