package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wantedly/risu/schema"
)

func TestCheckoutGitRepository(t *testing.T) {
	opts := schema.BuildCreateOpts{
		SourceRepo:   "wantedly/private-nginx-image-server",
		SourceBranch: "patched-small-light",
		ImageName:    "quay.io/wantedly/private-nginx-image-server:test",
	}
	build := schema.NewBuild(opts)
	err := checkoutGitRepository(build, "/tmp/risu/src/github.com/")
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat("/tmp/risu/src/github.com/wantedly/private-nginx-image-server/.git")
	if err != nil {
		t.Errorf("Fail to clone git repository\nerror: %v", err)
	}

	// Check for second try to test existing repository case
	err = checkoutGitRepository(build, "/tmp/risu/src/github.com/")
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat("/tmp/risu/src/github.com/wantedly/private-nginx-image-server/.git")
	if err != nil {
		t.Errorf("Fail to fetch&checkout git repository\nerror: %v", err)
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

func TestBuildFlow(t *testing.T) {
	response := httptest.NewRecorder()

	n := setUpServer()

	requestParams := `{
		"source_repo": "wantedly/risu",
		"source_revision": "ada9ce1829fab49e605e5a563dbf91274f64e923",
		"name": "quay.io/wantedly/risu:latest",
		"dockerfile": "Dockerfile.dev"
	}`

	// Create
	req, err := http.NewRequest("POST", "http://localhost:8080/builds", bytes.NewBuffer([]byte(requestParams)))
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("Got error for GET ruquest to /")
	}
	body := string(response.Body.Bytes())
	if body != "hoge" {
		t.Errorf("Got: %v", body)
	}

	// dec := json.NewDecoder(response.Body)
	// var build schema.Build
	// dec.Decode(&build)
	// expectedBuild := schema.Build{
	// 	ID:             uuid.NewUUID(),
	// 	SourceRepo:     "opts.SourceRepo",
	// 	SourceRevision: "opts.SourceRevision",
	// 	Name:           "opts.Name",
	// 	Dockerfile:     "Dockerfile",
	// 	Status:         "building",
	// 	CreatedAt:      time.Now(),
	// 	UpdatedAt:      time.Now(),
	// }
	//
	// if !build.Equals(expectedBuild) {
	// 	t.Errorf("Got empty body for GET request to /\n Got: %v\nExpected: %v", build, expectedBuild)
	// }
}
