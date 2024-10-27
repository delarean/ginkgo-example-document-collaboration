package collaboration_tests

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "collaboration/proto"
)

var _ = Describe("Check document's request:", func() {
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

	When("User requests a document", func() {
		It("must return current state of the document", func() {
			req := &pb.DocumentRequest{DocumentId: "doc1"}
			resp, err := client.GetDocument(ctx, req)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.DocumentId).To(Equal("doc1"))
			Expect(resp.Content).NotTo(BeNil())
			Expect(resp.Users).NotTo(BeEmpty())
		})
	})
})
