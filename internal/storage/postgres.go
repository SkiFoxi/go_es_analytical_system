package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akozadaev/go_es_analytical_system/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) Close() error {
	return ps.db.Close()
}

// GetBusinessTypes возвращает все типы бизнеса
func (ps *PostgresStorage) GetBusinessTypes(ctx context.Context) ([]*models.BusinessType, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM business_types ORDER BY name`

	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query business types: %w", err)
	}
	defer rows.Close()

	var businessTypes []*models.BusinessType
	for rows.Next() {
		var bt models.BusinessType
		if err := rows.Scan(
			&bt.ID,
			&bt.Name,
			&bt.Description,
			&bt.CreatedAt,
			&bt.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan business type: %w", err)
		}
		businessTypes = append(businessTypes, &bt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return businessTypes, nil
}

// GetRegions возвращает все регионы
func (ps *PostgresStorage) GetRegions(ctx context.Context) ([]*models.Region, error) {
	query := `SELECT id, name, parent_region_id, created_at, updated_at FROM regions ORDER BY name`

	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query regions: %w", err)
	}
	defer rows.Close()

	var regions []*models.Region
	for rows.Next() {
		var r models.Region
		var parentID sql.NullInt64
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&parentID,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan region: %w", err)
		}
		if parentID.Valid {
			parentIDInt := int(parentID.Int64)
			r.ParentRegionID = &parentIDInt
		}
		regions = append(regions, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return regions, nil
}
