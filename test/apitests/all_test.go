package apitests

import (
	"fmt"
	"net/http"
	"testing"
)

var initCatId string

func init() {
	// Preparation: delete all existing & create a cat
	ids := []string{}
	call("GET", "/cats", nil, nil, &ids)

	for _, id := range ids {
		code := 0
		call("DELETE", "/cats/"+id, nil, &code, nil)
		fmt.Println("DELETE /cats ->", code)
	}

	// Create a single cat into the DB
	call("POST", "/cats", &CatModel{Name: "Toto"}, nil, &initCatId)
}

// ---------- GET /cats ----------

func TestGetCats(t *testing.T) {

	code := 0
	result := []string{}
	err := call("GET", "/cats", nil, &code, &result)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats ->", code, result)

	if code != http.StatusOK {
		t.Error("We should get code 200, got", code)
	}

	// Après init(), la base ne contient qu'un seul chat
	if len(result) != 1 {
		t.Error("We should get one item, got", len(result))
		return
	}

	if result[0] != initCatId {
		t.Error("Listing the IDs, got", result[0])
	}
}

// ---------- POST /cats ----------

func TestCreateCat(t *testing.T) {
	code := 0
	newID := ""
	err := call("POST", "/cats", &CatModel{Name: "Rex", Color: "Black"}, &code, &newID)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("POST /cats ->", code, newID)

	if code != http.StatusCreated {
		t.Error("We should get code 201, got", code)
	}
	if newID == "" {
		t.Error("We should get the new cat ID")
	}
}

func TestCreateCatInvalidBody(t *testing.T) {
	// Cas limite : on envoie une simple chaîne au lieu d'un objet chat
	code := 0
	err := call("POST", "/cats", "pas un chat", &code, nil)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("POST /cats (invalid) ->", code)

	if code != http.StatusBadRequest {
		t.Error("We should get code 400, got", code)
	}
}

// ---------- GET /cats/{catId} ----------

func TestGetOneCat(t *testing.T) {
	code := 0
	cat := CatModel{}
	err := call("GET", "/cats/"+initCatId, nil, &code, &cat)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("GET /cats/{id} ->", code, cat)

	if code != http.StatusOK {
		t.Error("We should get code 200, got", code)
	}
	if cat.Name != "Toto" {
		t.Error("We should get the cat 'Toto', got", cat.Name)
	}
}

func TestGetOneCatNotFound(t *testing.T) {
	code := 0
	err := call("GET", "/cats/id-inconnu", nil, &code, nil)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("GET /cats/id-inconnu ->", code)

	if code != http.StatusNotFound {
		t.Error("We should get code 404, got", code)
	}
}

// ---------- DELETE /cats/{catId} ----------

func TestDeleteCat(t *testing.T) {
	// Préparation : on crée un chat dédié à ce test
	catID := ""
	call("POST", "/cats", &CatModel{Name: "ToDelete"}, nil, &catID)

	// Appel testé
	code := 0
	err := call("DELETE", "/cats/"+catID, nil, &code, nil)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("DELETE /cats/{id} ->", code)

	if code != http.StatusNoContent {
		t.Error("We should get code 204, got", code)
	}

	// Vérification : le chat ne doit plus exister
	getCode := 0
	call("GET", "/cats/"+catID, nil, &getCode, nil)
	if getCode != http.StatusNotFound {
		t.Error("The deleted cat should return 404, got", getCode)
	}
}

func TestDeleteCatNotFound(t *testing.T) {
	code := 0
	err := call("DELETE", "/cats/id-inconnu", nil, &code, nil)
	if err != nil {
		t.Fatal("Request error", err)
	}

	fmt.Println("DELETE /cats/id-inconnu ->", code)

	if code != http.StatusNotFound {
		t.Error("We should get code 404, got", code)
	}
}
