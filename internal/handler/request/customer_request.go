package request

type PathID struct {
	ID string `uri:"id" binding:"required,len=36"`
}

type CreateCustomerRequest struct {
	Name  string  `json:"name" validate:"required,max=255"`
	Email *string `json:"email" validate:"omitempty,email,max=255"`
	Phone *string `json:"phone" validate:"omitempty,phone,max=15"`
}

type UpdateCustomerRequest struct {
	Name  string  `json:"name" validate:"omitempty,max=255"`
	Email *string `json:"email" validate:"omitempty,email,max=255"`
	Phone *string `json:"phone" validate:"omitempty,phone,max=15"`
}
