package store

import (
	"fmt"
	"sync"

	"github.com/munizr13/bookshelf-api/internal/models"
)

type BookStore struct {
	mu    sync.RWMutex
	books map[string]models.Book
	nextID int
}

func New() *BookStore {
	return &BookStore{
		books:  make(map[string]models.Book),
		nextID: 1,
	}
}

func (s *BookStore) Add(b models.Book) models.Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.ID = fmt.Sprintf("book-%d", s.nextID)
	s.nextID++
	s.books[b.ID] = b
	return b
}

func (s *BookStore) Get(id string) (models.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *BookStore) List() []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

func (s *BookStore) Update(id string, b models.Book) (models.Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.books[id]
	if !ok {
		return models.Book{}, false
	}
	if b.Title != "" {
		existing.Title = b.Title
	}
	if b.Author != "" {
		existing.Author = b.Author
	}
	if b.Year != 0 {
		existing.Year = b.Year
	}
	existing.Read = b.Read
	s.books[id] = existing
	return existing, true
}

func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}
