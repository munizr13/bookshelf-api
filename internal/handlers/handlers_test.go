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

func TestUpdateBook(t *testing.T) {
	s := store.New()
	h := New(s)

	// Add a book first
	book := s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})

	body := `{"title":"Dune Messiah","year":1969}`
	req := httptest.NewRequest("PUT", "/books/"+book.ID, bytes.NewBufferString(body))
	req.SetPathValue("id", book.ID)
	w := httptest.NewRecorder()

	h.UpdateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var updated models.Book
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Title != "Dune Messiah" {
		t.Errorf("expected Dune Messiah, got %s", updated.Title)
	}
	if updated.Author != "Frank Herbert" {
		t.Errorf("expected author preserved, got %s", updated.Author)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	s := store.New()
	h := New(s)

	body := `{"title":"X"}`
	req := httptest.NewRequest("PUT", "/books/nonexistent", bytes.NewBufferString(body))
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	h.UpdateBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateBookInvalidJSON(t *testing.T) {
	s := store.New()
	h := New(s)

	req := httptest.NewRequest("PUT", "/books/any", bytes.NewBufferString("{bad"))
	req.SetPathValue("id", "any")
	w := httptest.NewRecorder()

	h.UpdateBook(w, req)

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
