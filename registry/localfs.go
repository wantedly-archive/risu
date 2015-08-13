package registry

import (
	"encoding/json"
	"os"

	"code.google.com/p/go-uuid/uuid"
	"github.com/wantedly/risu/schema"
)

// DefaultFilePath : default file path
const DefaultFilePath = "/tmp/risu/"

// LocalFsRegistry : sharing path
type LocalFsRegistry struct {
	path string
}

// NewLocalFsRegistry : init
func NewLocalFsRegistry(path string) Registry {
	if path == "" {
		path = DefaultFilePath
	}

	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0755)
	}
	return &LocalFsRegistry{path}
}

// Set : Stores the build data to a json file. file name is "/tmp/risu/<UUID>.json".
func (r *LocalFsRegistry) Set(build schema.Build) error {
	b, err := json.Marshal(build)
	if err != nil {
		return err
	}

	file, err := os.Create(DefaultFilePath + build.ID.String() + ".json")
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

// Get : get build data
func (r *LocalFsRegistry) Get(id uuid.UUID) (schema.Build, error) {
	file, err := os.Open(DefaultFilePath + id.String() + ".json")
	if err != nil {
		return schema.Build{}, err
	}
	defer file.Close()

	var build schema.Build

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&build)
	if err != nil {
		return schema.Build{}, err
	}

	return build, nil
}
