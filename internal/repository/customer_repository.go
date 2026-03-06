package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) domain.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindAll(limit, offset int, keyword string) ([]*domain.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers`
	var args []interface{}
	if keyword != "" {
		query += ` WHERE name LIKE ?`
		args = append(args, "%"+keyword+"%")
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer := &domain.Customer{}
		if err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *customerRepository) Count(keyword string) (int, error) {
	query := `SELECT COUNT(*) FROM customers`
	var args []interface{}

	if keyword != "" {
		query += ` WHERE name LIKE ?`
		args = append(args, "%"+keyword+"%")
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *customerRepository) FindByID(id string) (*domain.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers WHERE id = ? LIMIT 1`
	customer := &domain.Customer{}
	err := r.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or custom Not Found Error
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) Create(customer *domain.Customer) error {
	if customer.ID == "" {
		customer.ID = uuid.NewString()
	}
	query := `INSERT INTO customers (id, name, email, phone) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, customer.ID, customer.Name, customer.Email, customer.Phone)
	return err
}

func (r *customerRepository) Update(customer *domain.Customer) error {
	query := `UPDATE customers SET name = ?, email = ?, phone = ? WHERE id = ?`
	_, err := r.db.Exec(query, customer.Name, customer.Email, customer.Phone, customer.ID)
	return err
}

func (r *customerRepository) Delete(id string) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
