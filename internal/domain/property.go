package domain

import (
	"time"
)

type Property struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Rent      int        `json:"rent"`
	Address   string     `json:"address"`
	Layout    *string    `json:"layout"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PropertyRepository interface {
	FindAll() ([]*Property, error)
	FindByID(id string) (*Property, error)
	Create(property *Property) error
	Update(property *Property) error
	Delete(id string) error
}
