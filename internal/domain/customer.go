package domain

import (
	"time"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerListResult struct {
	TotalCount int         `json:"total_count"`
	Data       []*Customer `json:"data"`
}

type CustomerRepository interface {
	FindAll(limit, offset int, keyword string) ([]*Customer, error)
	Count(keyword string) (int, error)
	FindByID(id string) (*Customer, error)
	Create(customer *Customer) error
	Update(customer *Customer) error
	Delete(id string) error
}
