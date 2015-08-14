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

// NewBuild creates new build struct
func NewBuild(opts BuildCreateOpts) Build {
	if opts.SourceBranch == "" {
		opts.SourceBranch = "master"
	}

	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}

	currentTime := time.Now()
	return Build{
		ID:           uuid.NewUUID(),
		SourceRepo:   opts.SourceRepo,
		SourceBranch: opts.SourceBranch,
		Name:         opts.Name,
		Dockerfile:   opts.Dockerfile,
		Status:       "building",
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
}
