package main

import (
	"code.google.com/p/go-uuid/uuid"

	"github.com/wantedly/risu/registry"
	"github.com/wantedly/risu/schema"
)

func main() {
	reg := registry.NewRegistry("etcd", "http://172.17.8.101:4001")

	build := schema.Build{
		ID:             uuid.NewUUID(),
		SourceRepo:     "wantedly/risu",
		SourceRevision: "2c004f60b47bac66a3a83ffe40b822629251c037",
		Name:           "quay.io/wantedly/risu:latest",
		Dockerfile:     "Dockerfile",
	}

	err := reg.Set(build)
	if err != nil {
		panic(err)
	}

	_, err = reg.Get(build.ID)
	if err != nil {
		panic(err)
	}
}
