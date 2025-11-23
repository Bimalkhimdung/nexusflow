package repository

import (
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
)

// Repository implements data access layer
type Repository struct {
	db  *database.MultiTenantDB
	log *logger.Logger
}

// NewRepository creates a new repository instance
func NewRepository(db *database.DB, log *logger.Logger) *Repository {
	return &Repository{
		db:  database.NewMultiTenant(db),
		log: log,
	}
}

// Example data model
// type Model struct {
// 	database.BaseModel
// 	Name string `bun:"name,notnull"`
// }

// Example repository method
// func (r *Repository) GetByID(ctx context.Context, id string) (*Model, error) {
// 	var model Model
// 	err := r.db.NewSelect().
// 		Model(&model).
// 		Where("id = ?", id).
// 		Scan(ctx)
// 	
// 	if err != nil {
// 		r.log.Error("Failed to get model", logger.Error(err), logger.String("id", id))
// 		return nil, err
// 	}
// 	
// 	return &model, nil
// }

// func (r *Repository) Create(ctx context.Context, model *Model) error {
// 	_, err := r.db.NewInsert().
// 		Model(model).
// 		Exec(ctx)
// 	
// 	if err != nil {
// 		r.log.Error("Failed to create model", logger.Error(err))
// 		return err
// 	}
// 	
// 	return nil
// }
