package service

import (
	"context"
	"sync"

	pb "github.com/boseabhishek/foo-accounts/gen/go/accounts/proto/account"
)

// AccountServer implements the AccountService
type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	mu      sync.Mutex
	balance *pb.Money
}

// NewAccountServer creates a new AccountServer
func NewAccountServer() *AccountServer {
	return &AccountServer{
		balance: &pb.Money{Amount: 0},
	}
}

// Add implements the Add RPC
func (s *AccountServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.balance.Amount += req.Money.Amount

	return &pb.AddResponse{}, nil
}

// Show implements the Show RPC
func (s *AccountServer) Show(ctx context.Context, req *pb.ShowRequest) (*pb.ShowResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.ShowResponse{Money: s.balance}, nil
}
