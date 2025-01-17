package wasmtest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/request"
	"github.com/suborbital/reactr/rt"
)

func TestWasmRunnerWithFetchSwift(t *testing.T) {
	job := rt.NewJob("fetch-swift", "https://1password.com")

	res, err := sharedRT.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte))[:100] != "<!doctype html><html lang=en data-language-url=/><head><meta charset=utf-8><meta name=viewport conte" {
		t.Error(fmt.Errorf("expected 1password.com HTML, got %q", string(res.([]byte))[:100]))
	}
}

func TestWasmRunnerEchoSwift(t *testing.T) {
	job := rt.NewJob("hello-swift", "Connor")

	res, err := sharedRT.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello Connor" {
		t.Error(fmt.Errorf("hello Connor, got %s", string(res.([]byte))))
	}
}

func TestWasmRunnerSwift(t *testing.T) {
	body := testBody{
		Username: "cohix",
	}

	bodyJSON, _ := json.Marshal(body)

	req := &request.CoordinatedRequest{
		Method: "GET",
		URL:    "/hello/world",
		ID:     uuid.New().String(),
		Body:   bodyJSON,
		State: map[string][]byte{
			"hello": []byte("what is up"),
		},
	}

	reqJSON, err := req.ToJSON()
	if err != nil {
		t.Error("failed to ToJSON", err)
	}

	job := rt.NewJob("swift-log", reqJSON)

	res, err := sharedRT.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := &request.CoordinatedResponse{}
	if err := json.Unmarshal(res.([]byte), resp); err != nil {
		t.Error("failed to Unmarshal response")
	}

	if string(resp.Output) != "hello what is up" {
		t.Error(fmt.Errorf("expected 'hello, what is up', got %s", string(res.([]byte))))
	}
}

func TestWasmFileGetStaticSwift(t *testing.T) {
	getJob := rt.NewJob("get-static-swift", "")

	r, err := sharedRT.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Do get-static job"))
		return
	}

	result := string(r.([]byte))

	expected := "# Hello, World\n\nContents are very important"

	if result != expected {
		t.Error("failed, got:\n", result, "\nexpeted:\n", expected)
	}
}
