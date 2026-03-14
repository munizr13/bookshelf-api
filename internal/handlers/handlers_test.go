package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/munizr13/bookshelf-api/internal/models"
	"github.com/munizr13/bookshelf-api/internal/store"
)

func TestAddBook(t *testing.T) {
	s := store.New()
	h := New(s)

	body := `{"title":"Dune","author":"Frank Herbert","year":1965}`
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.AddBook(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var book models.Book
	json.NewDecoder(w.Body).Decode(&book)
	if book.Title != "Dune" {
		t.Errorf("expected Dune, got %s", book.Title)
	}
	if book.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestAddBookMissingTitle(t *testing.T) {
	s := store.New()
	h := New(s)

	body := `{"author":"Nobody"}`
	req := httptest.NewRequest("POST", "/books", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.AddBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetBookNotFound(t *testing.T) {
	s := store.New()
	h := New(s)

	req := httptest.NewRequest("GET", "/books/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	h.GetBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
