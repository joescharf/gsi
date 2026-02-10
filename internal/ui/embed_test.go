package ui

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDistFS(t *testing.T) {
	sub, err := DistFS()
	if err != nil {
		t.Fatalf("DistFS() error: %v", err)
	}

	f, err := sub.Open("index.html")
	if err != nil {
		t.Fatalf("Open index.html: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	if !strings.Contains(string(data), "<!DOCTYPE html>") {
		t.Error("index.html does not contain DOCTYPE")
	}
}

func TestHandler_Root(t *testing.T) {
	h, err := Handler()
	if err != nil {
		t.Fatalf("Handler() error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET / status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("GET / response does not contain DOCTYPE")
	}
}

func TestHandler_SPAFallback(t *testing.T) {
	h, err := Handler()
	if err != nil {
		t.Fatalf("Handler() error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboard/settings", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /dashboard/settings status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("SPA fallback response does not contain DOCTYPE")
	}
}

func TestHandler_MissingAsset404(t *testing.T) {
	h, err := Handler()
	if err != nil {
		t.Fatalf("Handler() error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/missing.js", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("GET /missing.js status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
