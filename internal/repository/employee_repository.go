package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	if employee.ID == "" {
		employee.ID = uuid.NewString()
	}
	query := `INSERT INTO employees (id, name, email, password_hash, role) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, employee.ID, employee.Name, employee.Email, employee.PasswordHash, employee.Role)
	return err
}

func (r *employeeRepository) Update(employee *domain.Employee) error {
	query := `UPDATE employees SET name = ?, email = ?, password_hash = ?, role = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, employee.Name, employee.Email, employee.PasswordHash, employee.Role, employee.ID)
	return err
}

func (r *employeeRepository) Delete(id string) error {
	query := `UPDATE employees SET deleted_at=NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *employeeRepository) FindAll(limit, offset int) ([]*domain.Employee, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM employees WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		employee := &domain.Employee{}
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Role,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (r *employeeRepository) FindByID(id string) (*domain.Employee, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM employees WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	employee := &domain.Employee{}
	err := r.db.QueryRow(query, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Email,
		&employee.Role,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return employee, nil
}
