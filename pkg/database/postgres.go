package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Config holds database configuration
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DB wraps bun.DB with additional functionality
type DB struct {
	*bun.DB
	config Config
}

// New creates a new database connection
func New(cfg Config) (*DB, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	// Parse config
	pgxConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Create connection pool
	sqldb := stdlib.OpenDB(*pgxConfig)

	// Set connection pool settings
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Create bun DB
	db := bun.NewDB(sqldb, pgdialect.New())

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:     db,
		config: cfg,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping checks database connectivity
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// BeginTx starts a new transaction
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

// RunInTx runs a function in a transaction
func (db *DB) RunInTx(ctx context.Context, opts *sql.TxOptions, fn func(ctx context.Context, tx bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// MultiTenantDB provides multi-tenant database operations
type MultiTenantDB struct {
	*DB
}

// NewMultiTenant creates a multi-tenant database wrapper
func NewMultiTenant(db *DB) *MultiTenantDB {
	return &MultiTenantDB{DB: db}
}

// WithOrgID returns a query scoped to an organization
func (db *MultiTenantDB) WithOrgID(orgID string) *bun.SelectQuery {
	return db.NewSelect().Where("organization_id = ?", orgID)
}

// ScopedInsert creates an insert query with organization_id
func (db *MultiTenantDB) ScopedInsert(ctx context.Context, orgID string, model interface{}) *bun.InsertQuery {
	return db.NewInsert().Model(model).Value("organization_id", "?", orgID)
}

// ScopedUpdate creates an update query scoped to organization
func (db *MultiTenantDB) ScopedUpdate(ctx context.Context, orgID string, model interface{}) *bun.UpdateQuery {
	return db.NewUpdate().Model(model).Where("organization_id = ?", orgID)
}

// ScopedDelete creates a delete query scoped to organization
func (db *MultiTenantDB) ScopedDelete(ctx context.Context, orgID string, model interface{}) *bun.DeleteQuery {
	return db.NewDelete().Model(model).Where("organization_id = ?", orgID)
}

// BaseModel provides common fields for all models
type BaseModel struct {
	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string    `bun:"organization_id,notnull"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	CreatedBy      string    `bun:"created_by"`
	UpdatedBy      string    `bun:"updated_by"`
	Version        int64     `bun:"version,notnull,default:1"`
}

// BeforeAppendModel hook for bun
func (m *BaseModel) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
		m.Version++
	}
	return nil
}
