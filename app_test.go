package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Vérifie que makeHandlerFunc renvoie bien le code et le JSON du service
func TestMakeHandlerFuncOK(t *testing.T) {
	fakeService := func(req *http.Request) (int, any) {
		return http.StatusOK, "hello"
	}
	handler := makeHandlerFunc(fakeService)

	// httptest.NewRecorder capture la réponse HTTP
	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest("GET", "/", nil))

	if recorder.Code != http.StatusOK {
		t.Error("On attend 200, obtenu", recorder.Code)
	}
	if recorder.Header().Get("content-type") != "application/json" {
		t.Error("Le content-type doit être application/json")
	}
	if recorder.Body.String() != "\"hello\"\n" {
		t.Error("Corps inattendu :", recorder.Body.String())
	}
}

// Cas limite : le service plante (panic), le serveur doit répondre 500 sans s'arrêter
func TestMakeHandlerFuncPanic(t *testing.T) {
	panicService := func(req *http.Request) (int, any) {
		panic("boom")
	}
	handler := makeHandlerFunc(panicService)

	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest("GET", "/", nil))

	if recorder.Code != http.StatusInternalServerError {
		t.Error("On attend 500 après une panic, obtenu", recorder.Code)
	}
}

// Test des routes complètes via newApp (le routeur)
func TestRouterListCats(t *testing.T) {
	app := newApp()

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, httptest.NewRequest("GET", "/api/cats", nil))

	if recorder.Code != http.StatusOK {
		t.Error("On attend 200 sur GET /api/cats, obtenu", recorder.Code)
	}
}
