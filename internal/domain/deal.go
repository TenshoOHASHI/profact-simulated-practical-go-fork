package domain

import (
	"time"
)

type Deal struct {
	ID           string     `json:"id"`
	CustomerID   string     `json:"customer_id"`
	CustomerName string     `json:"customer_name,omitempty"` // JOINなどで取得した場合
	PropertyID   *string    `json:"property_id"`
	AssigneeID   *string    `json:"assignee_id"`
	Status       string     `json:"status"`
	MoveInDate   *time.Time `json:"move_in_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type DealRepository interface {
	FindAll() ([]*Deal, error)
	FindByID(id string) (*Deal, error)
	Create(deal *Deal) error
	Update(deal *Deal) error
	UpdateStatus(id, status string) error
	Delete(id string) error
}
