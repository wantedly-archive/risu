package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/libgit2/git2go"
	"github.com/unrolled/render"

	"github.com/wantedly/risu/registry"
	"github.com/wantedly/risu/schema"
)

const (
	GitHubHost           = "github.com/"
	DefaultCloneBasePath = "/tmp/risu/repository/"
)

var ren = render.New()

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

	reg := registry.NewRegistry("localfs", "")
	err = reg.Set(build)
	if err != nil {
		log.Fatal(err)
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"status": "internal server error"})
		return
	}
	ren.JSON(w, http.StatusCreated, build)

	// git clone
	err = gitClone(build)
	if err != nil {
		log.Fatal(err)
	}
}

func root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ren.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	reg := registry.NewRegistry("localfs", "")
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
	reg := registry.NewRegistry("localfs", "")
	build, err := reg.Get(uuid)
	if err != nil {
		ren.JSON(w, http.StatusNotFound, map[string]string{"status": "not found"})
		return
	}
	ren.JSON(w, http.StatusOK, build)
}

// Clone run "git clone <repository_URL>" and "git checkout rebision"
func gitClone(build schema.Build) error {
	basePath := DefaultCloneBasePath
	if _, err := os.Stat(basePath); err != nil {
		os.MkdirAll(basePath, 0755)
	}

	// htpps://<token>@github.com/<SourceRepo>.git
	cloneURL := "https://" + os.Getenv("GITHUB_ACCESS_TOKEN") + "@" + GitHubHost + build.SourceRepo + ".git"

	// debug
	fmt.Println(cloneURL)

	clonePath := basePath + build.SourceRepo

	// debug
	fmt.Println(clonePath)

	_, err := git.Clone(cloneURL, clonePath, &git.CloneOptions{})
	if err != nil {
		return err
	}
	return nil
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
	if os.Getenv("GITHUB_ACCESS_TOKEN") == "" {
		os.Exit(1)
	}
	n := setUpServer()
	n.Run(":8080")
}
