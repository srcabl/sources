package bootstrap

import (
	pb "github.com/srcabl/protos/sources"
	"github.com/srcabl/sources/internal/config"
	"github.com/srcabl/sources/internal/service"
)

type Bootstrap struct {
	config  *config.Environment
	service pb.SourceServiceServer
}

func New(cfg *config.Environment) (*Bootstrap, error) {

	srvc, err := service.New()
	if err != nil {
		return nil, err
	}

	return &Bootstrap{
		config:  cfg,
		service: srvc,
	}, nil
}
