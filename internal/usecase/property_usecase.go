package usecase

import (
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type PropertyUsecase interface {
	ListProperties() ([]*domain.Property, error)
	GetProperty(id string) (*domain.Property, error)
	CreateProperty(property *domain.Property) error
	UpdateProperty(property *domain.Property) error
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

func (u *propertyUsecase) UpdateProperty(property *domain.Property) error {
	return u.repo.Update(property)
}

func (u *propertyUsecase) DeleteProperty(id string) error {
	return u.repo.Delete(id)
}
