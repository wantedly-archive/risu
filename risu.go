package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"github.com/wantedly/risu/registry"
	"github.com/wantedly/risu/schema"
	"github.com/wantedly/risu/shell"
)

const (
	SourceBasePath = "/var/risu/src/github.com/"
	CacheBasePath  = "/var/risu/cache"
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

	build := schema.NewBuild(opts)
	err = reg.Set(build)
	if err != nil {
		log.Fatal(err)
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"status": "internal server error"})
		return
	}
	ren.JSON(w, http.StatusAccepted, build)

	go func() {
		if err := gitClone(build); err != nil {
			return
		}

		if err := dockerBuild(build); err != nil {
			return
		}

		go dockerPush(build)
		go pushCache(build)
	}()
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

// Clone run "git clone <repository_URL>" and "git checkout branch"
func gitClone(build schema.Build) error {
	if _, err := os.Stat(SourceBasePath); err != nil {
		os.MkdirAll(SourceBasePath, 0755)
	}

	// htpps://<token>@github.com/<SourceRepo>.git
	cloneURL := "https://" + os.Getenv("GITHUB_ACCESS_TOKEN") + "@github.com/" + build.SourceRepo + ".git"

	// debug
	fmt.Println(cloneURL)

	clonePath := SourceBasePath + build.SourceRepo

	// debug
	fmt.Println(clonePath)

	_, err := shell.Command("git", "clone", cloneURL, clonePath)
	if err != nil {
		return err
	}
	return nil
}

func dockerBuild(build schema.Build) error {
	// TODO (@dtan4)
	return nil
}

func dockerPush(build schema.Build) error {
	// TODO (@koudaii)
	return nil
}

func pushCache(build schema.Build) error {
	// TODO (@dtan4)
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
		log.Fatal("Please provide 'GITHUB_ACCESS_TOKEN' through environment")
	}
	n := setUpServer()
	n.Run(":8080")
}
