package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wantedly/risu/schema"
)

func TestCheckoutGitRepository(t *testing.T) {
	opts := schema.BuildCreateOpts{
		SourceRepo: "wantedly/risu",
		ImageName:  "quay.io/wantedly/risu:test",
	}
	build := schema.NewBuild(opts)
	err := checkoutGitRepository(build, "/tmp/risu/src/github.com/")
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat("/tmp/risu/src/github.com/wantedly/risu/.git")
	if err != nil {
		t.Errorf("Fail to clone git repository\nerror: %v", err)
	}

	// Check for second try to test existing repository case
	err = checkoutGitRepository(build, "/tmp/risu/src/github.com/")
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat("/tmp/risu/src/github.com/wantedly/risu/.git")
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
		"source_branch": "ada9ce1829fab49e605e5a563dbf91274f64e923",
		"image_name": "quay.io/wantedly/risu:latest"
	}`

	// Create
	req, err := http.NewRequest("POST", "http://localhost:8080/builds", bytes.NewBuffer([]byte(requestParams)))
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Errorf("Got error for POST ruquest to /builds")
	}

	dec := json.NewDecoder(response.Body)
	var build schema.Build
	dec.Decode(&build)

	if build.SourceRepo != "wantedly/risu" ||
		build.SourceBranch != "ada9ce1829fab49e605e5a563dbf91274f64e923" ||
		build.ImageName != "quay.io/wantedly/risu:latest" ||
		build.Dockerfile != "Dockerfile" {
		t.Errorf("Create build failed.\nGot: %v", build)
	}

	uuid := build.ID.String()

	// Show
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:8080/builds/"+uuid, nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("Got error for Get ruquest to /builds/" + uuid)
	}

	dec = json.NewDecoder(response.Body)
	dec.Decode(&build)

	if build.ID.String() != uuid ||
		build.SourceRepo != "wantedly/risu" ||
		build.SourceBranch != "ada9ce1829fab49e605e5a563dbf91274f64e923" ||
		build.ImageName != "quay.io/wantedly/risu:latest" ||
		build.Dockerfile != "Dockerfile" {
		t.Errorf("Show build failed.\nGot: %v", build)
	}

	// Index
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:8080/builds", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("Got error for Get ruquest to /builds")
	}

	builds := make([]schema.Build, 0)
	dec = json.NewDecoder(response.Body)
	dec.Decode(&builds)
	if len(builds) == 0 {
		t.Errorf("Fail to index builds.\nGot: %v", builds)
	}
}
