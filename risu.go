package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"github.com/wantedly/risu/registry"
	"github.com/wantedly/risu/schema"
)

var ren = render.New()
var reg = registry.NewRegistry(os.Getenv("REGISTRY_BACKEND"), os.Getenv("REGISTRY_ENDPOINT"))

func create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	var opts schema.BuildCreateOpts
	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		log.Fatal(err)
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"status": "internal server error"})
		return
	}

	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}

	currentTime := time.Now()
	build := schema.Build{
		ID:             uuid.NewUUID(),
		SourceRepo:     opts.SourceRepo,
		SourceRevision: opts.SourceRevision,
		Name:           opts.Name,
		Dockerfile:     opts.Dockerfile,
		Status:         "building",
		CreatedAt:      currentTime,
		UpdatedAt:      currentTime,
	}

	err = reg.Set(build)
	if err != nil {
		log.Fatal(err)
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"status": "internal server error"})
		return
	}
	ren.JSON(w, http.StatusCreated, build)
}

func root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ren.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	builds, err := reg.List()
	if err != nil {
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"status": "internal server error"})
		return
	}
	ren.JSON(w, http.StatusOK, builds)
}

func show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	uuid := uuid.Parse(id)
	build, err := reg.Get(uuid)
	if err != nil {
		ren.JSON(w, http.StatusNotFound, map[string]string{"status": "not found"})
		return
	}
	ren.JSON(w, http.StatusOK, build)
}

func setUpServer() *negroni.Negroni {
	router := httprouter.New()
	router.GET("/", root)
	router.GET("/builds", index)
	router.GET("/builds/:id", show)
	router.POST("/builds", create)

	n := negroni.Classic()
	n.UseHandler(router)
	return n
}

func main() {
	n := setUpServer()
	n.Run(":8080")
}
