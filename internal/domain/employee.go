package domain

import "time"

type Employee struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password_hash string    `json:"-"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmployeesRepository interface {
	FindAll(limit, offset int) ([]*Employee, error)
	FindByID(id string) (*Employee, error)
	Create(employees *Employee) error
	Update(employees *Employee) error
	Delete(id string) error
}
