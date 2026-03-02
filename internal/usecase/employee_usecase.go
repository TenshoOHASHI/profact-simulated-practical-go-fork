package usecase

import (
	"errors"
	"strings"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type CreateEmployeeInput struct {
	Name     string
	Email    string
	Password string `json:"-"`
	Role     string
}

type UpdateEmployeeInput struct {
	ID       string
	Name     *string
	Email    *string
	Password *string `json:"-"`
	Role     *string
}

var (
	ErrDuplicateEmail = errors.New("このメールアドレスは既に登録されています")
)

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "UNIQUE constraint")
}

type EmployeeUsecase interface {
	ListEmployees(limit, offset int) ([]*domain.Employee, error)
	GetEmployee(id string) (*domain.Employee, error)
	CreateEmployee(employee *CreateEmployeeInput) (*domain.Employee, error)
	UpdateEmployee(input *UpdateEmployeeInput) (*domain.Employee, error)
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

func (u *employeeUsecase) CreateEmployee(input *CreateEmployeeInput) (*domain.Employee, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	employee := &domain.Employee{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
	}
	if err := u.repo.Create(employee); err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return employee, nil
}

func (u *employeeUsecase) UpdateEmployee(input *UpdateEmployeeInput) (*domain.Employee, error) {
	existing, err := u.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("ユーザーが見つかりません")
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}

	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		existing.PasswordHash = string(hashedPassword)
	}

	if input.Role != nil {
		existing.Role = *input.Role
	}

	if err := u.repo.Update(existing); err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return existing, nil
}

func (u *employeeUsecase) DeleteEmployee(id string) error {
	return u.repo.Delete(id)
}
