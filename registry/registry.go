package registry

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/wantedly/risu/schema"
)

type Registry interface {
	Set(build schema.Build) error
	Get(id uuid.UUID) (schema.Build, error)
}

func NewRegistry(backend string, endpoint string) Registry {
	switch backend {
	case "etcd":
		return NewEtcdRegistry(endpoint)
	default:
		return NewEtcdRegistry(endpoint)
	}
}
