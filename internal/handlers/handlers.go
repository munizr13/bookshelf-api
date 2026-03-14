package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/munizr13/bookshelf-api/internal/models"
	"github.com/munizr13/bookshelf-api/internal/store"
)

type Handler struct {
	store *store.BookStore
}

func New(s *store.BookStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books := h.store.List()
	writeJSON(w, http.StatusOK, books)
}

func (h *Handler) AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if book.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	created := h.store.Add(book)
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	book, ok := h.store.Get(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	updated, ok := h.store.Update(id, book)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !h.store.Delete(id) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
