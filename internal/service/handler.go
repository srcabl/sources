package service

import (
	"context"

	"github.com/pkg/errors"
	sharedpb "github.com/srcabl/protos/shared"
	pb "github.com/srcabl/protos/sources"
	"github.com/srcabl/services/pkg/db/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler implements the sources service
type Handler struct {
	pb.UnimplementedSourcesServiceServer
	srcDeterminer SourceDeterminer
	datarepo      DataRepository
}

// New creates the service handler
func New(db *mysql.Client) (*Handler, error) {
	dataRepo, err := NewDataRepository(db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create data repo")
	}
	srcDeterminer, err := NewSourceDeterminer(dataRepo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create source determiner")
	}
	return &Handler{
		srcDeterminer: srcDeterminer,
		datarepo:      dataRepo,
	}, nil
}

// HealthCheck is the base healthcheck for the service
func (h *Handler) HealthCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {

	return nil, nil
}

// DetermineLinkSource determines the source of a post
func (h *Handler) DetermineLinkSource(ctx context.Context, req *pb.DetermineLinkSourceRequest) (*pb.DetermineLinkSourceResponse, error) {
	sources, err := h.srcDeterminer.DetermineSource(req.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrapf(err, "failed to deteremine the source of %s", req.Url).Error())
	}
	var dbSources []*DBSource
	for _, s := range sources {
		dbs, err := HydrateSourceModelForCreateFromSource(s.Source)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, errors.Wrapf(err, "failed to hydrate source %+v", s).Error())
		}
		dbSources = append(dbSources, dbs)
	}
	if err := h.datarepo.CreateSource(ctx, dbSources); err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "something happened creating sources").Error())
	}
	var pbSources []*sharedpb.SourceNode
	for _, s := range dbSources {
		pbs, err := s.ToGRPC()
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, errors.Wrapf(err, "failed to hydrate source %+v", s).Error())
		}
		pbSources = append(pbSources, &sharedpb.SourceNode{Source: pbs})
	}
	return &pb.DetermineLinkSourceResponse{
		PrimarySourceNodes: pbSources,
	}, nil
}
