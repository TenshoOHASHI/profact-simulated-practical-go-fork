package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type dealRepository struct {
	db *sql.DB
}

func NewDealRepository(db *sql.DB) domain.DealRepository {
	return &dealRepository{db: db}
}

func (r *dealRepository) FindAll() ([]*domain.Deal, error) {
	query := `
		SELECT
			d.id, d.customer_id, d.property_id, d.assignee_id, d.status, d.move_in_date, d.created_at, d.updated_at,
			c.name as customer_name,
			p.name as property_name,
			e.name as assignee_name
		FROM deals d
		JOIN customers c ON d.customer_id = c.id AND c.deleted_at IS NULL
		LEFT JOIN properties p ON d.property_id = p.id AND p.deleted_at IS NULL
  	    LEFT JOIN employees e ON d.assignee_id = e.id AND e.deleted_at IS NULL
		WHERE d.deleted_at IS NULL
		ORDER BY d.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deals []*domain.Deal
	for rows.Next() {
		deal := &domain.Deal{}
		if err := rows.Scan(
			&deal.ID,
			&deal.CustomerID,
			&deal.PropertyID,
			&deal.AssigneeID,
			&deal.Status,
			&deal.MoveInDate,
			&deal.CreatedAt,
			&deal.UpdatedAt,
			&deal.CustomerName,
			&deal.PropertyName,
			&deal.AssigneeName,
		); err != nil {
			return nil, err
		}
		deals = append(deals, deal)
	}
	return deals, nil
}

func (r *dealRepository) FindByID(id string) (*domain.Deal, error) {
	query := `
		SELECT
			d.id, d.customer_id, d.property_id, d.assignee_id, d.status, d.move_in_date, d.created_at, d.updated_at,
			c.name as customer_name,
			p.name as property_name,
			e.name as assignee_name
		FROM deals d
		JOIN customers c ON d.customer_id = c.id AND c.deleted_at IS NULL
		LEFT JOIN properties p ON d.property_id = p.id AND p.deleted_at IS NULL
  		LEFT JOIN employees e ON d.assignee_id = e.id AND e.deleted_at IS NULL
		WHERE d.id = ? AND d.deleted_at IS NULL LIMIT 1`

	deal := &domain.Deal{}
	err := r.db.QueryRow(query, id).Scan(
		&deal.ID,
		&deal.CustomerID,
		&deal.PropertyID,
		&deal.AssigneeID,
		&deal.Status,
		&deal.MoveInDate,
		&deal.CreatedAt,
		&deal.UpdatedAt,
		&deal.CustomerName,
		&deal.PropertyName,
		&deal.AssigneeName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return deal, nil
}

func (r *dealRepository) Create(deal *domain.Deal) error {
	if deal.ID == "" {
		deal.ID = uuid.NewString()
	}
	query := `INSERT INTO deals (id, customer_id, property_id, assignee_id, status, move_in_date) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, deal.ID, deal.CustomerID, deal.PropertyID, deal.AssigneeID, deal.Status, deal.MoveInDate)
	return err
}

func (r *dealRepository) Update(deal *domain.Deal) error {
	query := `UPDATE deals SET customer_id = ?, property_id = ?, assignee_id = ?, status = ?, move_in_date = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, deal.CustomerID, deal.PropertyID, deal.AssigneeID, deal.Status, deal.MoveInDate, deal.ID)
	return err
}

func (r *dealRepository) UpdateStatus(id, status string) error {
	query := `UPDATE deals SET status = ?, assignee_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *dealRepository) Delete(id string) error {
	query := `UPDATE deals SET deleted_at= NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}
