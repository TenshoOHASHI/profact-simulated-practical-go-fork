package usecase

import (
	"errors"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

// 簡易的なステータスバリデーション例
var validStatuses = map[string]bool{
	"new_lead":          true,
	"following_up":      true,
	"viewing_scheduled": true,
	"applying":          true,
	"contracted":        true,
	"lost":              true,
}

type DealUsecase interface {
	ListDeals() ([]*domain.Deal, error)
	GetDeal(id string) (*domain.Deal, error)
	CreateDeal(deal *domain.Deal) error
	UpdateDeal(deal *domain.Deal) (*domain.Deal, error)
	UpdateDealStatus(id, status string) (*domain.Deal, error)
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

func (u *dealUsecase) UpdateDeal(deal *domain.Deal) (*domain.Deal, error) {
	existing, err := u.repo.FindByID(deal.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("deal not found")
	}

	if deal.CustomerName != "" {
		existing.CustomerName = deal.CustomerName
	}
	if deal.PropertyID != nil {
		existing.PropertyID = deal.PropertyID
	}

	if deal.AssigneeID != nil {
		existing.AssigneeID = deal.AssigneeID
	}

	if deal.Status != "" && !validStatuses[deal.Status] {
		return nil, errors.New("invalid status transition")
	}

	if deal.Status != "" {
		existing.Status = deal.Status
	}

	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, err
}

func (u *dealUsecase) UpdateDealStatus(id, status string) (*domain.Deal, error) {

	existing, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("deal not found")
	}

	if status != "" && !validStatuses[status] {
		return nil, errors.New("invalid status transition")
	}

	// ステータスが指定されている場合のみ更新
	if status != "" {
		existing.Status = status
	}

	if err := u.repo.UpdateStatus(id, existing.Status); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *dealUsecase) DeleteDeal(id string) error {
	return u.repo.Delete(id)
}
