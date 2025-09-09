package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// un faux ServiceFunc pour tester
func dummyServiceFunc(req *http.Request) (int, any) {
	return http.StatusOK, map[string]string{"message": "ok"}
}

// un faux ServiceFunc qui panic pour tester la recovery
func panicServiceFunc(req *http.Request) (int, any) {
	panic("boom!")
}

func TestMakeHandlerFunc_Success(t *testing.T) {
	// Arrange
	handler := makeHandlerFunc(dummyServiceFunc)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Vérifie que la réponse est bien du JSON attendu
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body["message"] != "ok" {
		t.Errorf("expected message 'ok', got '%s'", body["message"])
	}
}

func TestMakeHandlerFunc_Panic(t *testing.T) {
	// Arrange
	handler := makeHandlerFunc(panicServiceFunc)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var body string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("expected body '%s', got '%s'", http.StatusText(http.StatusInternalServerError), body)
	}
}

func TestLogReq(t *testing.T) {
	// Arrange : un handler factice qui écrit "hello"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // 418
	})
	handler := logReq(finalHandler)

	req := httptest.NewRequest(http.MethodGet, "/log-test", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status %d, got %d", http.StatusTeapot, rec.Code)
	}
}

func TestNewApp_RoutesExist(t *testing.T) {
	app := newApp()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"POST", "/api/cats"},
		{"GET", "/api/cats"},
		{"GET", "/api/cats/123"},
		{"DELETE", "/api/cats/123"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		// Ici tu ne connais pas la réponse exacte, mais tu peux vérifier que le routeur
		// ne renvoie pas un 404
		if rec.Code == http.StatusNotFound {
			t.Errorf("route %s %s not found", tc.method, tc.path)
		}
	}
}
