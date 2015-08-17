package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

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

func NewEtcdRegistry(machines string) Registry {
	if os.Getenv("RISU_ETCD_MACHINES") != "" {
		machines = os.Getenv("RISU_ETCD_MACHINES")
	}

	if machines == "" {
		machines = DefaultMachines
	}

	m := strings.Split(machines, ",")
	etcdClient := *etcd.NewClient(m)
	return &EtcdRegistry{etcdClient, DefaultKeyPrefix}
}

func (r *EtcdRegistry) Create(opts schema.BuildCreateOpts) (schema.Build, error) {
	build := schema.NewBuild(&opts)
	j, err := marshal(build)
	if err != nil {
		return build, err
	}

	key := path.Join(r.keyPrefix, build.ID.String())
	_, err = r.etcd.Create(key, string(j), 0)
	if err != nil {
		return build, err
	}

	return build, nil
}

func (r *EtcdRegistry) Set(build schema.Build) error {
	j, err := marshal(build)
	if err != nil {
		return err
	}

	key := path.Join(r.keyPrefix, build.ID.String())
	_, err = r.etcd.Set(key, string(j), 0)
	if err != nil {
		return err
	}

	return nil
}

func (r *EtcdRegistry) Get(id uuid.UUID) (schema.Build, error) {
	key := path.Join(r.keyPrefix, id.String())
	res, err := r.etcd.Get(key, false, true)
	if err != nil {
		return schema.Build{}, err
	}

	var build schema.Build
	err = unmarshal(res.Node.Value, &build)
	if err != nil {
		if isKeyNotFound(err) {
			// TODO (spesnova): return 404 error
			return schema.Build{}, err
		}
		return schema.Build{}, err
	}

	return build, nil
}

func (r *EtcdRegistry) List() ([]schema.Build, error) {
	var builds []schema.Build

	key := path.Join(r.keyPrefix)
	res, err := r.etcd.Get(key, false, true)
	if err != nil {
		if isKeyNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, node := range res.Node.Nodes {
		var build schema.Build
		err = unmarshal(node.Value, &build)
		if err != nil {
			return nil, err
		}

		builds = append(builds, build)
	}

	return builds, nil
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
