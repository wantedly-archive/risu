package registry

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/wantedly/risu/schema"
)

type Registry interface {
	Set(id uuid.UUID, build schema.Build) error
	Get(id uuid.UUID) (schema.Build, error)
}
