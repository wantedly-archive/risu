package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wantedly/risu/schema"
)

func TestGitClone(t *testing.T) {
	opts := schema.BuildCreateOpts{
		SourceRepo: "wantedly/private-nginx-image-server",
		Name:       "quay.io/wantedly/private-nginx-image-server:test",
	}
	build := schema.NewBuild(opts)
	err := gitClone(build)
	if err != nil {
		t.Error(err)
	}
}

func TestRootAccess(t *testing.T) {
	response := httptest.NewRecorder()

	n := setUpServer()

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("Got error for GET ruquest to /")
	}
	body := string(response.Body.Bytes())
	expectedBody := "{\"status\":\"ok\"}"
	if body != expectedBody {
		t.Errorf("Got empty body for GET request to /\n Got: %s, Expected: %s", body, expectedBody)
	}
}
