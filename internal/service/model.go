package service

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	pb "github.com/srcabl/protos/shared"
	sourcepb "github.com/srcabl/protos/sources"
	"github.com/srcabl/services/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DBSource is the database source model
type DBSource struct {
	UUID          string
	Name          string
	Origanization string
	CreatedByUUID string
	CreatedAt     int64
	UpdatedByUUID sql.NullString
	UpdatedAt     sql.NullInt64
}

// CreatedByUUIDString satisfies the services helper to transform db auditfields to grpc auditfields
func (s *DBSource) CreatedByUUIDString() string {
	return s.CreatedByUUID
}

// CreatedAtUnixInt satisfies the services helper to transform db auditfields to grpc auditfields
func (s *DBSource) CreatedAtUnixInt() int64 {
	return s.CreatedAt
}

// UpdatedByUUIDNullString satisfies the services helper to transform db auditfields to grpc auditfields
func (s *DBSource) UpdatedByUUIDNullString() sql.NullString {
	return s.UpdatedByUUID
}

// UpdatedAtUnixNullInt satisfies the services helper to transform db auditfields to grpc auditfields
func (s *DBSource) UpdatedAtUnixNullInt() sql.NullInt64 {
	return s.UpdatedAt
}

// ToGRPC transforms the dbuser to proto link
func (s *DBSource) ToGRPC() (*pb.Source, error) {
	id, err := uuid.FromString(s.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to transform uuid: %s", s.UUID)
	}
	auditFields, err := proto.DBAuditFieldsToGRPC(s)
	if err != nil {
		return nil, errors.Wrap(err, "failed to transform auditfields")
	}
	return &pb.Source{
		Uuid:         id.Bytes(),
		Name:         s.Name,
		Organization: s.Origanization,
		AuditFields:  auditFields,
	}, nil
}

// HydrateSourceModelForCreateFromSource creates a dbsource from a proto source
func HydrateSourceModelForCreateFromSource(src *pb.Source) (*DBSource, error) {
	source, err := hydrateSource(src.Name, src.Organization)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hydrate sourced from proto source %+v", src)
	}
	return source, nil
}

// HydrateSourceModelForCreateFromRequest creates a db post from a proto post and fills in any missing data
func HydrateSourceModelForCreateFromRequest(req *sourcepb.CreateSourceRequest) (*DBSource, error) {
	source, err := hydrateSource(req.Name, req.Organization)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hydrate sourced from create request %+v", req)
	}
	return source, nil
}

func hydrateSource(name, org string) (*DBSource, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "failed to generate uuid for link").Error())
	}
	now := time.Now().Unix()
	return &DBSource{
		UUID:          newUUID.String(),
		Name:          name,
		Origanization: org,
		CreatedByUUID: newUUID.String(),
		CreatedAt:     now,
		UpdatedByUUID: sql.NullString{Valid: true, String: newUUID.String()},
		UpdatedAt:     sql.NullInt64{Valid: true, Int64: now},
	}, nil

}
