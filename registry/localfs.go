package registry

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"code.google.com/p/go-uuid/uuid"

	"github.com/wantedly/risu/schema"
)

// DefaultFilePath is default file path
const DefaultFilePath = "/tmp/risu/"

// LocalFsRegistry is sharing path
type LocalFsRegistry struct {
	path string
}

// NewLocalFsRegistry is init
func NewLocalFsRegistry(path string) Registry {
	if path == "" {
		path = DefaultFilePath
	}

	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0755)
	}
	return &LocalFsRegistry{path}
}

// Set stores the build data to a json file. file name is "/tmp/risu/<UUID>.json".
func (r *LocalFsRegistry) Set(build schema.Build) error {
	b, err := json.Marshal(build)
	if err != nil {
		return err
	}

	file, err := os.Create(r.path + build.ID.String() + ".json")
	if err != nil {
		return err
	}

	defer file.Close()

	buildData := []byte(string(b))

	_, err = file.Write(buildData)
	if err != nil {
		return err
	}
	return nil
}

// Get get build data
func (r *LocalFsRegistry) Get(id uuid.UUID) (schema.Build, error) {
	// TODO: add error handling.(Expired Data and Not Found File)
	file, err := os.Open(r.path + id.String() + ".json")
	if err != nil {
		return schema.Build{}, err
	}
	defer file.Close()

	var build schema.Build

	err = json.NewDecoder(file).Decode(&build)
	if err != nil {
		return schema.Build{}, err
	}

	return build, nil
}

func (r *LocalFsRegistry) List() ([]schema.Build, error) {
	var builds []schema.Build

	fileInfos, err := ioutil.ReadDir(r.path)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		file, err := os.Open(r.path + fileInfo.Name())
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var build schema.Build

		err = json.NewDecoder(file).Decode(&build)
		if err != nil {
			return nil, err
		}

		builds = append(builds, build)
	}

	return builds, nil
}
