package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

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
	router.GET("/build/:image", show)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")
}
