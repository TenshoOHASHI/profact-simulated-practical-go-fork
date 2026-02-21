package usecase

import (
	"errors"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type DealUsecase interface {
	ListDeals() ([]*domain.Deal, error)
	GetDeal(id string) (*domain.Deal, error)
	CreateDeal(deal *domain.Deal) error
	UpdateDeal(deal *domain.Deal) error
	UpdateDealStatus(id, status string, assigneeID *string) error
	DeleteDeal(id string) error
}

type dealUsecase struct {
	repo domain.DealRepository
}

func NewDealUsecase(repo domain.DealRepository) DealUsecase {
	return &dealUsecase{repo: repo}
}

func (u *dealUsecase) ListDeals() ([]*domain.Deal, error) {
	return u.repo.FindAll()
}

func (u *dealUsecase) GetDeal(id string) (*domain.Deal, error) {
	return u.repo.FindByID(id)
}

func (u *dealUsecase) CreateDeal(deal *domain.Deal) error {
	return u.repo.Create(deal)
}

func (u *dealUsecase) UpdateDeal(deal *domain.Deal) error {
	return u.repo.Update(deal)
}

func (u *dealUsecase) UpdateDealStatus(id, status string, assigneeID *string) error {
	// 簡易的なステータスバリデーション例
	validStatuses := map[string]bool{
		"new_lead":          true,
		"following_up":      true,
		"viewing_scheduled": true,
		"application":       true,
		"contract":          true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status transition")
	}

	return u.repo.UpdateStatus(id, status, assigneeID)
}

func (u *dealUsecase) DeleteDeal(id string) error {
	return u.repo.Delete(id)
}
