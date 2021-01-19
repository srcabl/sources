package service

import (
	"context"
	"errors"

	pb "github.com/srcabl/protos/sources"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler implments the sources service
type Handler struct {
	pb.UnimplementedSourceServiceServer
}

// New creates the service handler
func New() (*Handler, error) {
	return &Handler{}, nil
}

// HealthCheck is the base healthcheck for the service
func (h *Handler) HealthCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {

	return nil, nil
}

// DeterminePostSource determines the source of a post
func (h *Handler) DeterminePostSource(ctx context.Context, req *pb.DeterminePostSourceRequest) (*pb.DeterminePostSourceResponse, error) {

	return nil, errors.New("not implemented")
}
