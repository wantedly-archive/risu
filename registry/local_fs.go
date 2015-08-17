package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/wantedly/risu/schema"
)

// DefaultFileDir is default file dir
const DefaultFileDir = "/etc/risu/"

// LocalFsRegistry is sharing dir
type LocalFsRegistry struct {
	dir string
}

// NewLocalFsRegistry is init
func NewLocalFsRegistry(dir string) Registry {
	if dir == "" {
		dir = DefaultFileDir
	}

	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	go func(dir string) {
		for {
			if err := cleanUpItems(dir); err != nil {
				log.Fatal(err)
			}

			time.Sleep(10 * time.Second)
		}
	}(dir)

	return &LocalFsRegistry{dir}
}

func (r *LocalFsRegistry) Create(opts schema.BuildCreateOpts) (schema.Build, error) {
	build := schema.NewBuild(&opts)
	b, err := json.Marshal(build)
	if err != nil {
		return build, err
	}

	file, err := os.Create(r.dir + build.ID.String() + ".json")
	if err != nil {
		return build, err
	}

	defer file.Close()

	buildData := []byte(string(b))

	_, err = file.Write(buildData)
	if err != nil {
		return build, err
	}
	return build, nil
}

// Set stores the build data to a json file. file name is "/tmp/risu/<UUID>.json".
func (r *LocalFsRegistry) Set(build schema.Build, opts schema.BuildUpdateOpts) error {
	build = schema.UpdateBuild(build, &opts)
	b, err := json.Marshal(build)
	if err != nil {
		return err
	}

	file, err := os.Create(r.dir + build.ID.String() + ".json")
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
	file, err := os.Open(r.dir + id.String() + ".json")
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

	fileInfos, err := ioutil.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		file, err := os.Open(r.dir + fileInfo.Name())
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

func deleteExpiredItem(path string, f os.FileInfo, _ error) error {
	if f.IsDir() {
		return nil
	}

	elapsedSeconds := time.Now().Sub(f.ModTime()).Seconds()

	if elapsedSeconds > DefaultExpireSeconds {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

func cleanUpItems(dir string) error {
	if err := filepath.Walk(dir, deleteExpiredItem); err != nil {
		return err
	}

	return nil
}
