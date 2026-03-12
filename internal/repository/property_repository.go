package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type propertyRepository struct {
	db *sql.DB
}

func NewPropertyRepository(db *sql.DB) domain.PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) FindAll() ([]*domain.Property, error) {
	query := `SELECT id, name, rent, address, layout, status, created_at, updated_at FROM properties WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []*domain.Property
	for rows.Next() {
		property := &domain.Property{}
		if err := rows.Scan(
			&property.ID,
			&property.Name,
			&property.Rent,
			&property.Address,
			&property.Layout,
			&property.Status,
			&property.CreatedAt,
			&property.UpdatedAt,
		); err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}
	return properties, nil
}

func (r *propertyRepository) FindByID(id string) (*domain.Property, error) {
	query := `SELECT id, name, rent, address, layout, status, created_at, updated_at FROM properties WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	property := &domain.Property{}
	err := r.db.QueryRow(query, id).Scan(
		&property.ID,
		&property.Name,
		&property.Rent,
		&property.Address,
		&property.Layout,
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or custom Not Found Error
		}
		return nil, err
	}
	return property, nil
}

func (r *propertyRepository) Create(property *domain.Property) error {
	if property.ID == "" {
		property.ID = uuid.NewString()
	}
	query := `INSERT INTO properties (id, name, rent, address, layout, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, property.ID, property.Name, property.Rent, property.Address, property.Layout, property.Status)
	return err
}

func (r *propertyRepository) Update(property *domain.Property) error {
	query := `UPDATE properties SET name = ?, rent = ?, address = ?, layout = ?, status = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, property.Name, property.Rent, property.Address, property.Layout, property.Status, property.ID)
	return err
}

func (r *propertyRepository) Delete(id string) error {
	query := `UPDATE properties SET deleted_at=NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *propertyRepository) GetExistingPropertiesMap() (map[string]bool, error) {
	query := `SELECT name ,address FROM properties WHERE deleted_at IS NULL`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var name, address string
		if err := rows.Scan(&name, &address); err != nil {
			return nil, err
		}
		key := name + "|" + address
		existing[key] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *propertyRepository) BulkCreate(properties []*domain.Property) error {
	if len(properties) == 0 {
		return nil
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	valueStrings := make([]string, 0, len(properties))
	valueArgs := make([]interface{}, 0, len(properties))

	for _, p := range properties {
		if p.ID == "" {
			p.ID = uuid.NewString()
		}
		valueStrings = append(valueStrings, "(?,?,?,?,?,?)")
		valueArgs = append(valueArgs, p.ID, p.Name, p.Rent, p.Address, p.Layout, p.Status)

	}
	query := fmt.Sprintf(
		"INSERT INTO properties (id, name, rent, address, layout, status) VALUES %s",
		strings.Join(valueStrings, ","),
	)
	_, err = tx.Exec(query, valueArgs...)

	if err != nil {
		return err
	}
	return tx.Commit()
}
