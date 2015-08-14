package schema

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Build struct {
	ID           uuid.UUID `json:"id"`
	SourceRepo   string    `json:"source_repo"`
	SourceBranch string    `json:"source_branch"`
	Name         string    `json:"name"`
	Dockerfile   string    `json:"dockerfile"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BuildCreateOpts struct {
	SourceRepo   string `json:"source_repo"`
	SourceBranch string `json:"source_branch"`
	Name         string `json:"name"`
	Dockerfile   string `json:"dockerfile"`
}
