package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/srcabl/services/pkg/db/mysql"
)

// DataRepositoryGetter defines the behavior of the data repo getters
type DataRepositoryGetter interface {
	GetSource(context.Context, string) error
}

// DataRepositoryCreator defines the behavior of the data repo creators
type DataRepositoryCreator interface {
	CreateSource(context.Context, []*DBSource) error
}

// DataRepository defines the behavior of the data repo
type DataRepository interface {
	DataRepositoryGetter
	DataRepositoryCreator
}

type dataRepository struct {
	db *mysql.Client
}

// NewDataRepository news up a sources data repo
func NewDataRepository(db *mysql.Client) (DataRepository, error) {
	return &dataRepository{
		db: db,
	}, nil
}

// GetSorce gets a source
func (dr *dataRepository) GetSource(ctx context.Context, uuid string) error {
	//TODO: this......all of this
	return nil
}

const createSourceStatement = `
INSERT INTO
	sources(
		uuid,
		name,
		organization,
		created_by_uuid,
		created_at,
		updated_by_uuid,
		updated_at
	)
VALUES
	(?, ?, ?, ?, ?, ?, ?)
`

// CreateSource creates a source
func (dr *dataRepository) CreateSource(ctx context.Context, sources []*DBSource) error {
	tx, err := dr.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to begin transaction")
	}
	stm, err := tx.PrepareContext(ctx, createSourceStatement)
	if err != nil {
		if rollErr := tx.Rollback(); rollErr != nil {
			return errors.Wrapf(rollErr, "failed to rollback after failing to create source %+v", sources)
		}
		return errors.Wrapf(err, "failed to prepare statement to create source %+v", sources)
	}
	for _, source := range sources {
		_, err = stm.ExecContext(ctx,
			source.UUID,
			source.Name,
			source.Origanization,
			source.CreatedByUUID,
			source.CreatedAt,
			source.UpdatedByUUID.String,
			source.UpdatedAt.Int64,
		)
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				return errors.Wrapf(rollErr, "failed to rollback after failing to create source %+v", source)
			}
			return errors.Wrapf(err, "failed to execute statment to create source %+v", source)
		}
	}
	if err := tx.Commit(); err != nil {
		if rollErr := tx.Rollback(); rollErr != nil {
			return errors.Wrapf(rollErr, "failed to rollback after failing to create source %+v", sources)
		}
		return errors.Wrapf(err, "failed to create source %+v", sources)
	}
	return nil
}
