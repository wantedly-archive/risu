package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"github.com/wantedly/risu/registry"
	"github.com/wantedly/risu/schema"
)

var ren = render.New()

func create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	var opts schema.BuildCreateOpts
	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}

	build := schema.Build{
		ID:             uuid.NewUUID(),
		SourceRepo:     opts.SourceRepo,
		SourceRevision: opts.SourceRevision,
		Name:           opts.Name,
		Dockerfile:     opts.Dockerfile,
		Status:         "building",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	reg := registry.NewRegistry("localfs", "")
	reg.Set(build)

	// debug code
	builddata, err := reg.Get(build.ID)
	fmt.Fprintln(w, builddata)
}

func root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ren.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	fmt.Fprintf(w, "Welcome, %s!\n", name)
}

func show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	uuid := uuid.Parse(id)
	reg := registry.NewRegistry("localfs", "")
	build, err := reg.Get(uuid)
	if err != nil {
		ren.JSON(w, http.StatusNotFound, map[string]string{"status": "not found"})
	}
	ren.JSON(w, http.StatusOK, build)
}

func main() {
	router := httprouter.New()
	router.GET("/", root)
	router.GET("/builds", index)
	router.GET("/builds/:id", show)
	router.POST("/builds", create)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")
}
