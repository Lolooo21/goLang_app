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

	if len(result) != 1 {
		t.Error("We should get one item, got", len(result))
		return
	}

	if result[0] != initCatId {
		t.Error("Listing the IDs, got", result[0])
	}
}

func TestGetCat(t *testing.T) {
	code := 0
	result := CatModel{}

	// Appel avec l'ID qu'on a créé dans init()
	err := call("GET", "/cats/"+initCatId, nil, &code, &result)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats/{id} ->", code, result)

	if code != http.StatusOK {
		t.Error("Expected 200, got", code)
	}

	if result.ID != initCatId {
		t.Errorf("Expected cat ID %s, got %s", initCatId, result.ID)
	}
	if result.Name != "Toto" {
		t.Errorf("Expected cat name 'Toto', got '%s'", result.Name)
	}
}

func TestGetCat_NotFound(t *testing.T) {
	code := 0
	message := ""

	// Appel avec un ID inexistant
	err := call("GET", "/cats/unknown-id", nil, &code, &message)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats/unknown-id ->", code, message)

	if code != http.StatusNotFound {
		t.Error("Expected 404, got", code)
	}
	if message != "Cat not found" {
		t.Errorf("Expected message 'Cat not found', got '%s'", message)
	}
}

func TestDeleteCat(t *testing.T) {
	code := 0

	// Suppression du chat qu'on a créé
	err := call("DELETE", "/cats/"+initCatId, nil, &code, nil)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("DELETE /cats/{id} ->", code)

	if code != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", code)
	}

	// Vérifier que le chat n'existe plus
	code = 0
	message := ""
	err = call("GET", "/cats/"+initCatId, nil, &code, &message)
	if err != nil {
		t.Error("Request error", err)
	}

	if code != http.StatusNotFound {
		t.Error("Expected 404 after deletion, got", code)
	}
}

func TestDeleteCat_NotFound(t *testing.T) {
	code := 0
	message := ""

	// Suppression d'un chat inexistant
	err := call("DELETE", "/cats/fake-id", nil, &code, &message)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("DELETE /cats/fake-id ->", code, message)

	if code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", code)
	}
	if message != "Cat not found" {
		t.Errorf("Expected message 'Cat not found', got '%s'", message)
	}
}
