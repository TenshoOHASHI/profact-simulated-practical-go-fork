package request

type CreateDealRequest struct {
	CustomerID string `json:"customer_id" validate:"required,len=36"`
	PropertyID string `json:"property_id" validate:"required,len=36"`
	AssigneeID string `json:"assignee_id" validate:"required,len=36"`
	Status     string `json:"status" validate:"required,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
	MoveInDate string `json:"move_in_date" validate:"required,datetime=2006-01-02T15:04:05Z"`
}

type UpdateDealRequest struct {
	CustomerID string  `json:"customer_id" validate:"omitempty,len=36"`
	PropertyID *string `json:"property_id" validate:"omitempty,len=36"`
	AssigneeID *string `json:"assignee_id" validate:"omitempty,len=36"`
	Status     string  `json:"status" validate:"omitempty,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
	MoveInDate *string `json:"move_in_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z"`
}

type UpdateDealStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=new_lead following_up viewing_scheduled applying contracted lost"`
}
