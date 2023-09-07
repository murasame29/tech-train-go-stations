package model

import "time"

type (
	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct {
		TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)

type TODO struct {
	ID          int64     `json:"id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTODORequest struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}
