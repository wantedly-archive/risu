package registry

import (
	"encoding/json"
	"io/ioutil"
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

	return &LocalFsRegistry{path}
}

// Set : build meta data save. file name is "/tmp/risu/UUID.json".
func (r *LocalFsRegistry) Set(build schema.Build) error {
	b, err := json.Marshal(build)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(DefaultFilePath+string(build.ID)+".json", []byte(string(b)), 0666)
	if err != nil {
		return err
	}
	return nil
}

// Get : get build
func (r *LocalFsRegistry) Get(id uuid.UUID) (schema.Build, error) {
	targetBuildData, err := os.Open(DefaultFilePath + string(id) + ".json")
	if err != nil {
		return schema.Build{}, err
	}
	defer targetBuildData.Close()

	var build schema.Build

	decoder := json.NewDecoder(targetBuildData)

	err = decoder.Decode(build)
	if err != nil {
		return schema.Build{}, err
	}

	return build, nil
}
