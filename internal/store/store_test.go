package store

import (
	"testing"

	"github.com/munizr13/bookshelf-api/internal/models"
)

func TestAddAndGet(t *testing.T) {
	s := New()
	book := s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})

	if book.ID == "" {
		t.Fatal("expected ID to be set")
	}

	got, ok := s.Get(book.ID)
	if !ok {
		t.Fatal("expected to find book")
	}
	if got.Title != "Dune" {
		t.Errorf("expected title Dune, got %s", got.Title)
	}
}

func TestDelete(t *testing.T) {
	s := New()
	book := s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949})

	if !s.Delete(book.ID) {
		t.Fatal("expected delete to succeed")
	}
	if _, ok := s.Get(book.ID); ok {
		t.Fatal("expected book to be gone")
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	book := s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})

	updated, ok := s.Update(book.ID, models.Book{Title: "Dune Messiah", Year: 1969})
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Title != "Dune Messiah" {
		t.Errorf("expected title Dune Messiah, got %s", updated.Title)
	}
	if updated.Author != "Frank Herbert" {
		t.Errorf("expected author preserved, got %s", updated.Author)
	}
	if updated.Year != 1969 {
		t.Errorf("expected year 1969, got %d", updated.Year)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := New()
	_, ok := s.Update("nonexistent", models.Book{Title: "X"})
	if ok {
		t.Fatal("expected update to fail for missing book")
	}
}

func TestList(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Book A"})
	s.Add(models.Book{Title: "Book B"})

	books := s.List()
	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}
}
