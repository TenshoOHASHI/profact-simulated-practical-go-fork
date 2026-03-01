package usecase

import (
	"errors"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeUsecase interface {
	ListEmployees(limit, offset int) ([]*domain.Employee, error)
	GetEmployee(id string) (*domain.Employee, error)
	CreateEmployee(employee *domain.Employee) error
	UpdateEmployee(employee *domain.Employee) (*domain.Employee, error)
	DeleteEmployee(id string) error
}

type employeeUsecase struct {
	repo domain.EmployeeRepository
}

func NewEmployeeUsecase(repo domain.EmployeeRepository) EmployeeUsecase {
	return &employeeUsecase{repo: repo}
}

func (u *employeeUsecase) ListEmployees(limit, offset int) ([]*domain.Employee, error) {
	if limit > 20 {
		limit = 20
	}
	return u.repo.FindAll(limit, offset)
}

func (u *employeeUsecase) GetEmployee(id string) (*domain.Employee, error) {
	return u.repo.FindByID(id)
}

func (u *employeeUsecase) CreateEmployee(employee *domain.Employee) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	employee.PasswordHash = string(hashedPassword)
	return u.repo.Create(employee)
}

func (u *employeeUsecase) UpdateEmployee(employee *domain.Employee) (*domain.Employee, error) {
	existing, err := u.repo.FindByID(employee.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	if employee.Name != "" {
		existing.Name = employee.Name
	}
	if employee.Email != "" {
		existing.Email = employee.Email
	}

	if employee.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		existing.PasswordHash = string(hashedPassword)
	}

	if employee.Role != "" {
		existing.Role = employee.Role
	}

	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *employeeUsecase) DeleteEmployee(id string) error {
	return u.repo.Delete(id)
}
