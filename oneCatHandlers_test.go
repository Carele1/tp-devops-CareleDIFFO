package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Petite fonction utilitaire : crée une requête avec le paramètre {catId} rempli
func newRequestWithID(method, catID string) *http.Request {
	req := httptest.NewRequest(method, "/api/cats/"+catID, nil)
	req.SetPathValue("catId", catID)
	return req
}

// ---------- getCat ----------

func TestGetCatOK(t *testing.T) {
	// Préparation : on place un chat connu dans la base
	catsDatabase["test-get"] = Cat{Name: "Garfield", Color: "Orange"}

	code, result := getCat(newRequestWithID("GET", "test-get"))

	if code != http.StatusOK {
		t.Error("On attend 200, obtenu", code)
	}
	cat, ok := result.(Cat)
	if !ok {
		t.Fatal("Le résultat devrait être un Cat")
	}
	if cat.Name != "Garfield" {
		t.Error("Nom attendu 'Garfield', obtenu", cat.Name)
	}
}

func TestGetCatNotFound(t *testing.T) {
	// Cas limite : ID qui n'existe pas
	code, _ := getCat(newRequestWithID("GET", "id-inconnu"))

	if code != http.StatusNotFound {
		t.Error("On attend 404, obtenu", code)
	}
}

// ---------- deleteCat ----------

func TestDeleteCatOK(t *testing.T) {
	catsDatabase["test-delete"] = Cat{Name: "Tom"}

	code, _ := deleteCat(newRequestWithID("DELETE", "test-delete"))

	if code != http.StatusNoContent {
		t.Error("On attend 204, obtenu", code)
	}
	// Vérification : le chat ne doit plus être dans la base
	if _, found := catsDatabase["test-delete"]; found {
		t.Error("Le chat aurait dû être supprimé")
	}
}

func TestDeleteCatNotFound(t *testing.T) {
	code, _ := deleteCat(newRequestWithID("DELETE", "id-inconnu"))

	if code != http.StatusNotFound {
		t.Error("On attend 404, obtenu", code)
	}
}
