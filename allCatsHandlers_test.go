package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ---------- listMapKeys ----------

func TestListMapKeysEmpty(t *testing.T) {
	// Préparation : une map vide
	input := map[string]Cat{}

	// Appel de la fonction testée
	result := listMapKeys(input)

	// Vérification : on attend une liste vide (et pas nil)
	if result == nil {
		t.Error("Le résultat ne doit pas être nil")
	}
	if len(result) != 0 {
		t.Error("On attend 0 clé, obtenu", len(result))
	}
}

func TestListMapKeysSeveral(t *testing.T) {
	input := map[string]Cat{
		"a": {Name: "Toto"},
		"b": {Name: "Felix"},
	}

	result := listMapKeys(input)

	if len(result) != 2 {
		t.Fatal("On attend 2 clés, obtenu", len(result))
	}
	// L'ordre d'une map n'est pas garanti en Go : on vérifie juste la présence
	for _, key := range []string{"a", "b"} {
		found := false
		for _, r := range result {
			if r == key {
				found = true
			}
		}
		if !found {
			t.Error("La clé est absente du résultat :", key)
		}
	}
}

// ---------- createCat ----------

func TestCreateCatOK(t *testing.T) {
	// Préparation : une requête HTTP simulée avec un JSON valide
	body := strings.NewReader(`{"name": "Rex", "color": "Black"}`)
	req := httptest.NewRequest("POST", "/api/cats", body)

	code, result := createCat(req)

	if code != http.StatusCreated {
		t.Error("On attend 201, obtenu", code)
	}

	// Le résultat doit être l'ID du chat, et le chat doit être dans la "base"
	newID, ok := result.(string)
	if !ok {
		t.Fatal("Le résultat devrait être un ID (string)")
	}
	cat, found := catsDatabase[newID]
	if !found {
		t.Fatal("Le chat créé n'est pas dans la base")
	}
	if cat.Name != "Rex" {
		t.Error("Nom attendu 'Rex', obtenu", cat.Name)
	}
}

func TestCreateCatInvalidJSON(t *testing.T) {
	// Cas limite : JSON cassé
	body := strings.NewReader(`{"name": `)
	req := httptest.NewRequest("POST", "/api/cats", body)

	code, _ := createCat(req)

	if code != http.StatusBadRequest {
		t.Error("On attend 400 pour un JSON invalide, obtenu", code)
	}
}
