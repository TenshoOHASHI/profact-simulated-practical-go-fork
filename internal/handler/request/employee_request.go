package request

type CreateEmployeeRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin manager agent"`
}
type UpdateEmployeeRequest struct {
	Name     *string `json:"name" validate:"omitempty,max=255"`
	Email    *string `json:"email" validate:"omitempty,email,max=255"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin manager agent"`
}
