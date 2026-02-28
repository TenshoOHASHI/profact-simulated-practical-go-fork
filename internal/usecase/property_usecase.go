package usecase

import (
	"errors"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

// 簡易的なステータスバリデーション例
var validPropertyStatuses = map[string]bool{
	"available":  true,
	"contracted": true,
	"hidden":     true,
}

type PropertyUsecase interface {
	ListProperties() ([]*domain.Property, error)
	GetProperty(id string) (*domain.Property, error)
	CreateProperty(property *domain.Property) error
	UpdateProperty(property *domain.Property) (*domain.Property, error)
	DeleteProperty(id string) error
}

type propertyUsecase struct {
	repo domain.PropertyRepository
}

func NewPropertyUsecase(repo domain.PropertyRepository) PropertyUsecase {
	return &propertyUsecase{repo: repo}
}

func (u *propertyUsecase) ListProperties() ([]*domain.Property, error) {
	return u.repo.FindAll()
}

func (u *propertyUsecase) GetProperty(id string) (*domain.Property, error) {
	return u.repo.FindByID(id)
}

func (u *propertyUsecase) CreateProperty(property *domain.Property) error {
	return u.repo.Create(property)
}

func (u *propertyUsecase) UpdateProperty(property *domain.Property) (*domain.Property, error) {
	existing, err := u.repo.FindByID(property.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("property not found")
	}

	if property.Name != "" {
		existing.Name = property.Name
	}

	if property.Rent != 0 {
		existing.Rent = property.Rent
	}

	if property.Address != "" {
		existing.Address = property.Address
	}

	if property.Layout != nil {
		existing.Layout = property.Layout
	}

	if property.Status != "" && !validPropertyStatuses[property.Status] {
		return nil, errors.New("invalid property status transition")
	}

	if property.Status != "" {
		existing.Status = property.Status
	}
	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, err
}

func (u *propertyUsecase) DeleteProperty(id string) error {
	return u.repo.Delete(id)
}
