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

func TestSearchByAuthor(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949})
	s.Add(models.Book{Title: "Children of Dune", Author: "Frank Herbert", Year: 1976})

	author := "herbert"
	books := s.Search(SearchFilters{Author: &author})
	if len(books) != 2 {
		t.Errorf("expected 2 books by Herbert, got %d", len(books))
	}
}

func TestSearchByRead(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Read: true})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Read: false})
	s.Add(models.Book{Title: "Neuromancer", Author: "William Gibson", Read: true})

	read := true
	books := s.Search(SearchFilters{Read: &read})
	if len(books) != 2 {
		t.Errorf("expected 2 read books, got %d", len(books))
	}

	unread := false
	books = s.Search(SearchFilters{Read: &unread})
	if len(books) != 1 {
		t.Errorf("expected 1 unread book, got %d", len(books))
	}
}

func TestSearchByYear(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949})
	s.Add(models.Book{Title: "Foundation", Author: "Isaac Asimov", Year: 1951})

	year := 1965
	books := s.Search(SearchFilters{Year: &year})
	if len(books) != 1 {
		t.Errorf("expected 1 book from 1965, got %d", len(books))
	}
	if books[0].Title != "Dune" {
		t.Errorf("expected Dune, got %s", books[0].Title)
	}
}

func TestSearchCombinedFilters(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Dune", Author: "Frank Herbert", Year: 1965, Read: true})
	s.Add(models.Book{Title: "Children of Dune", Author: "Frank Herbert", Year: 1976, Read: false})
	s.Add(models.Book{Title: "1984", Author: "George Orwell", Year: 1949, Read: true})

	author := "herbert"
	read := true
	books := s.Search(SearchFilters{Author: &author, Read: &read})
	if len(books) != 1 {
		t.Errorf("expected 1 read Herbert book, got %d", len(books))
	}
	if books[0].Title != "Dune" {
		t.Errorf("expected Dune, got %s", books[0].Title)
	}
}

func TestSearchNoFilters(t *testing.T) {
	s := New()
	s.Add(models.Book{Title: "Dune"})
	s.Add(models.Book{Title: "1984"})

	books := s.Search(SearchFilters{})
	if len(books) != 2 {
		t.Errorf("expected 2 books with no filters, got %d", len(books))
	}
}
