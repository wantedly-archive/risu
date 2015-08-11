package registry

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/go-etcd/etcd"

	"github.com/wantedly/risu/schema"
)

const (
	DefaultMachines  = "http://172.17.8.101:4001"
	DefaultKeyPrefix = "/risu/"
)

type EtcdRegistry struct {
	etcd      etcd.Client
	keyPrefix string
}

func NewRegistry(machines string, keyPrefix string) *Registry {
	if os.Getenv("RISU_ETCD_MACHINES") != "" {
		machines = os.Getenv("RISU_ETCD_MACHINES")
	}

	if machines == "" {
		machines = DefaultMachines
	}

	m := strings.Split(machines, ",")
	etcdClient := *etcd.NewClient(m)
	return &Registry{etcdClient, keyPrefix}
}

func (r *EtcdRegistry) Set(id uuid.UUID, build schema.Build) error {
	return nil
}

func (r *EtcdRegistry) Get(id uuid.UUID) (schema.Build, error) {
	return nil
}

func marshal(obj interface{}) (string, error) {
	encoded, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("unable to JSON-serialize object: %s", err)
	}
	return string(encoded), nil
}

func unmarshal(val string, obj interface{}) error {
	err := json.Unmarshal([]byte(val), &obj)
	if err != nil {
		return fmt.Errorf("unable to JSON-deserialize object: %s", err)
	}
	return nil
}

func isKeyNotFound(err error) bool {
	e, ok := err.(*etcd.EtcdError)
	return ok && e.ErrorCode == 100
}
