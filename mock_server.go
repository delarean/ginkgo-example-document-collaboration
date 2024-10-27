package collaboration

import (
	"context"
	"sync"
	"time"

	pb "collaboration/proto"
)

type mockCollaborationServer struct {
	pb.UnimplementedCollaborationServiceServer
	mu              sync.Mutex
	documentContent map[string]string
	subscribers     map[string][]pb.CollaborationService_SubscribeToDocumentServer
}

func NewMockCollaborationServer() *mockCollaborationServer {
	return &mockCollaborationServer{
		documentContent: make(map[string]string),
		subscribers:     make(map[string][]pb.CollaborationService_SubscribeToDocumentServer),
	}
}

func (s *mockCollaborationServer) SubscribeToDocument(req *pb.DocumentSubscriptionRequest, stream pb.CollaborationService_SubscribeToDocumentServer) error {
	s.mu.Lock()
	s.subscribers[req.DocumentId] = append(s.subscribers[req.DocumentId], stream)
	s.mu.Unlock()
	// Блокируем поток, чтобы он оставался открытым
	for {
		select {
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *mockCollaborationServer) SendEdit(ctx context.Context, req *pb.EditRequest) (*pb.EditResponse, error) {
	s.mu.Lock()
	s.mu.Unlock()
	// Обновляем содержимое документа
	content := s.documentContent[req.DocumentId]
	// Простая обработка вставки (для примера)
	if req.Edit.Type == pb.EditType_INSERT {
		content = content[:req.Edit.Position] + req.Edit.Text + content[req.Edit.Position:]
	}
	s.documentContent[req.DocumentId] = content

	// Рассылаем обновление подписчикам
	for _, subscriber := range s.subscribers[req.DocumentId] {
		update := &pb.DocumentUpdate{
			DocumentId: req.DocumentId,
			Edits:      []*pb.Edit{req.Edit},
		}
		sub := subscriber
		err := sub.Send(update)
		if err != nil {
			return nil, err
		}
	}

	return &pb.EditResponse{
		Success:   true,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (s *mockCollaborationServer) GetDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	content := s.documentContent[req.DocumentId]
	return &pb.DocumentResponse{
		DocumentId: req.DocumentId,
		Content:    content,
		Users: []*pb.User{
			{UserId: "user1", DisplayName: "User One"},
			{UserId: "user2", DisplayName: "User Two"},
		},
	}, nil
}
