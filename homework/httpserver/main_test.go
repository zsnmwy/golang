package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheckHandler(t *testing.T) {
	// make new request
	req, err := http.NewRequest("Get", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// record Response
	rr := httptest.NewRecorder()

	// define test handler
	handler := http.HandlerFunc(healthzHandler)

	// get response from fake request
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `ok`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHeaderCopy(t *testing.T) {

	expectedKey := time.Now().String()
	expectedValues := time.Now().String()

	// make new request
	req, err := http.NewRequest("Get", "/", nil)
	req.Header.Set(expectedKey, expectedValues)
	if err != nil {
		t.Fatal(err)
	}

	// record Response
	rr := httptest.NewRecorder()

	// define test handler
	handler := http.HandlerFunc(rootHandler)

	// get response from fake request
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the header is what we expect
	if rr.Header().Get(expectedKey) != expectedValues {
		t.Errorf("handler returned unexpected header key: got %v want %v",
			rr.Header().Get(expectedKey), expectedValues)
	}
}

func TestVersion(t *testing.T) {
	// make new request
	req, err := http.NewRequest("Get", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// record Response
	rr := httptest.NewRecorder()

	// define test handler
	handler := http.HandlerFunc(rootHandler)

	// get response from fake request
	handler.ServeHTTP(rr, req)

	// check version string
	fmt.Println(rr.Header().Get("VERSION"))
	if rr.Header().Get("VERSION") == "" {
		t.Errorf("handler returned wrong version: got nil")
	}
}
