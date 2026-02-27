package request

import "time"

type CreateDealRequest struct {
	CustomerID string     `json:"customer_id" validate:"required,len=36"`
	PropertyID string     `json:"property_id" validate:"required,len=36"`
	AssigneeID string     `json:"assignee_id" validate:"required,len=36"`
	Status     string     `json:"status" validate:"required,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
	MoveInDate *time.Time `json:"move_in_date" validate:"required"`
}

type UpdateDealRequest struct {
	CustomerID string     `json:"customer_id" validate:"omitempty,len=36"`
	PropertyID *string    `json:"property_id" validate:"omitempty,len=36"`
	AssigneeID *string    `json:"assignee_id" validate:"omitempty,len=36"`
	Status     string     `json:"status" validate:"omitempty,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
	MoveInDate *time.Time `json:"move_in_date"`
}

type UpdateDealStatusRequest struct {
	Status     string  `json:"status" validate:"required,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
	AssigneeID *string `json:"assignee_id" validate:"omitempty,len=36"`
}
