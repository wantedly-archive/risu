package registry

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/wantedly/risu/schema"
)

const (
	DefaultExpireSeconds = 60 * 60 * 24 * 5 // 5 days
)

type Registry interface {
	Create(opts schema.BuildCreateOpts) (schema.Build, error)
	Set(build schema.Build, opts schema.BuildUpdateOpts) error
	Get(id uuid.UUID) (schema.Build, error)
	List() ([]schema.Build, error)
}

func NewRegistry(backend string, endpoint string) Registry {
	switch backend {
	case "etcd":
		return NewEtcdRegistry(endpoint)
	case "localfs":
		return NewLocalFsRegistry(endpoint)
	default:
		return NewLocalFsRegistry(endpoint)
	}
}
