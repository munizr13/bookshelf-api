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

func TestListBooksFilterByAuthor(t *testing.T) {
	s := store.New()
	h := New(s)
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949})

	req := httptest.NewRequest("GET", "/books?author=herbert", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []models.Book
	json.NewDecoder(w.Body).Decode(&books)
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].Title != "Dune" {
		t.Errorf("expected Dune, got %s", books[0].Title)
	}
}

func TestListBooksFilterByRead(t *testing.T) {
	s := store.New()
	h := New(s)
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Read: true})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Read: false})

	req := httptest.NewRequest("GET", "/books?read=true", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []models.Book
	json.NewDecoder(w.Body).Decode(&books)
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].Title != "Dune" {
		t.Errorf("expected Dune, got %s", books[0].Title)
	}
}

func TestListBooksFilterByYear(t *testing.T) {
	s := store.New()
	h := New(s)
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949})

	req := httptest.NewRequest("GET", "/books?year=1949", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []models.Book
	json.NewDecoder(w.Body).Decode(&books)
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].Title != "1984" {
		t.Errorf("expected 1984, got %s", books[0].Title)
	}
}

func TestListBooksInvalidRead(t *testing.T) {
	s := store.New()
	h := New(s)

	req := httptest.NewRequest("GET", "/books?read=notabool", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListBooksInvalidYear(t *testing.T) {
	s := store.New()
	h := New(s)

	req := httptest.NewRequest("GET", "/books?year=abc", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListBooksCombinedFilters(t *testing.T) {
	s := store.New()
	h := New(s)
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965, Read: true})
	s.Add(models.Book{Title: "Children of Dune", Author: "Frank Herbert", Year: 1976, Read: false})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949, Read: true})

	req := httptest.NewRequest("GET", "/books?author=herbert&read=true", nil)
	w := httptest.NewRecorder()
	h.ListBooks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []models.Book
	json.NewDecoder(w.Body).Decode(&books)
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].Title != "Dune" {
		t.Errorf("expected Dune, got %s", books[0].Title)
	}
}
