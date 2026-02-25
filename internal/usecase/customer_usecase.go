package usecase

import (
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
)

type CustomerUsecase interface {
	ListCustomers(limit, offset int) ([]*domain.Customer, error)
	GetCustomer(id string) (*domain.Customer, error)
	CreateCustomer(customer *domain.Customer) error
	UpdateCustomer(customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(id string) error
}

type customerUsecase struct {
	repo domain.CustomerRepository
}

func NewCustomerUsecase(repo domain.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo: repo}
}

func (u *customerUsecase) ListCustomers(limit, offset int) ([]*domain.Customer, error) {
	return u.repo.FindAll(limit, offset)
}

func (u *customerUsecase) GetCustomer(id string) (*domain.Customer, error) {
	return u.repo.FindByID(id)
}

func (u *customerUsecase) CreateCustomer(customer *domain.Customer) error {
	// ここにバリデーションなどのビジネスロジックを必要に応じて追加
	return u.repo.Create(customer)
}

func (u *customerUsecase) UpdateCustomer(customer *domain.Customer) (*domain.Customer, error) {
	existing, err := u.repo.FindByID(customer.ID)
	if err != nil {
		return nil, err
	}

	if customer.Name != "" {
		existing.Name = customer.Name
	}
	if customer.Email != nil {
		existing.Email = customer.Email
	}

	if customer.Phone != nil {
		existing.Phone = customer.Phone
	}
	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, err
}

func (u *customerUsecase) DeleteCustomer(id string) error {
	return u.repo.Delete(id)
}
