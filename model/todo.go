package model

import "time"

type (
	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct {
		TODO `json:"todo"`
	}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct {
		TODOs []TODO `json:"todos"`
	}

	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct {
		TODO `json:"todo"`
	}

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

type UpdateTODORequest struct {
	ID          int64  `json:"id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type ReadTODORequest struct {
	PrevID int64 `json:"prev_id"`
	Size   int64 `json:"size"`
}
