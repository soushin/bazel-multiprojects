package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestServer(t *testing.T) {

	msg = "Awesome"

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler)

	testHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Awesome Go! - built by Bazel!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
