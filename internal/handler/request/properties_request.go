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

type CSVRow struct {
	LineNumber int
	Name       string
	Rent       int
	Address    string
	Layout     string
	Status     string
}

type ValidationError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ImportResult struct {
	ImportedCount int           `json:"imported_count"`
	SkippedCount  int           `json:"skipped_count"`
	SkippedItems  []SkippedItem `json:"skipped_details,omitempty"`
}

type SkippedItem struct {
	Row    int    `json:"row,omitempty"`
	Reason string `json:"reason"`
}
