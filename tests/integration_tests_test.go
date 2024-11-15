package tests

import (
	"RTDS_API/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartStreamHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/stream/start", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "your-secure-api-key")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.StartStreamHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
