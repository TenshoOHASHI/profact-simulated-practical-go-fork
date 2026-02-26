package request

type CreatePropertyRequest struct {
	Name    string `json:"name" validate:"required,max=255"`
	Rent    int    `json:"rent" validate:"required,min=1,max=100000000"`
	Address string `json:"address" validate:"required,max=255"`
	Layout  string `json:"layout" validate:"required,max=50"`
	Status  string `json:"status" validate:"required,oneof=available contracted hidden"`
}
type UpdatePropertyRequest struct {
	Name    string  `json:"name" validate:"omitempty,max=255"`
	Rent    int     `json:"rent" validate:"omitempty,min=1,max=100000000"`
	Address string  `json:"address" validate:"omitempty,max=255"`
	Layout  *string `json:"layout" validate:"omitempty,max=50"`
	Status  string  `json:"status" validate:"omitempty,oneof=available contracted hidden"`
}
