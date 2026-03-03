package domain

import "time"

type Employee struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmployeeRepository interface {
	FindAll(limit, offset int) ([]*Employee, error)
	FindByID(id string) (*Employee, error)
	Create(employee *Employee) error
	Update(employee *Employee) error
	Delete(id string) error
}
