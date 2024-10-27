package collaboration_tests

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "collaboration/proto"
)

var _ = Describe("Check document's subscription:", func() {
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

	When("user subscribed to the document", Focus, func() {
		It("must receive it's updates", func() {
			var stream pb.CollaborationService_SubscribeToDocumentClient
			var err error

			By("first client is subscribing to the document", func() {
				stream, err = client.SubscribeToDocument(ctx, &pb.DocumentSubscriptionRequest{
					DocumentId: "doc1",
					User: &pb.User{
						UserId:      "user1",
						DisplayName: "User One",
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			var wg sync.WaitGroup
			var update *pb.DocumentUpdate

			wg.Add(1)
			By("first client try to receive the update from the stream", func() {
				go func() {
					defer wg.Done()
					// we should add it to see panic message inside goroutine
					defer GinkgoRecover()

					update, err = stream.Recv()
					Expect(err).NotTo(HaveOccurred())
				}()
			})
			time.Sleep(10 * time.Millisecond)

			By("another client sends an edit", func() {
				anotherClient := pb.NewCollaborationServiceClient(conn)
				_, err = anotherClient.SendEdit(ctx, &pb.EditRequest{
					DocumentId: "doc1",
					User: &pb.User{
						UserId:      "user2",
						DisplayName: "User Two",
					},
					Edit: &pb.Edit{
						UserId:    "user2",
						Position:  0,
						Text:      "Update from user2",
						Type:      pb.EditType_INSERT,
						Timestamp: time.Now().Unix(),
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

			By("first client has received the update", func() {
				wg.Wait()
				Expect(update.DocumentId).To(Equal("doc1"))
				Expect(update.Edits).To(HaveLen(1))
				Expect(update.Edits[0].Text).To(Equal("Update from user2"))
			})
		})
	})
})
