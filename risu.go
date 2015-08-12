package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Build is "Create a new build.""
type Build struct {
	SourceRepo     string `json:"source_repo"`
	SourceRevision string `json:"source_revision"`
	Name           string `json:"name"`
	Dockerfile     string `json:"dockerfile"`
}

func create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	build := &Build{Dockerfile: "Dockerfile"} // default setup Dockerfile
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&build)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, build)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	fmt.Fprintf(w, "Welcome, %s!\n", name)
}

func show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	image := ps.ByName("image")
	fmt.Fprintf(w, "Build %s!\n", image)
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/builds/:image", show)
	router.POST("/builds", create)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")
}
