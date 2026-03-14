package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(nil)
		log.SetFlags(log.LstdFlags)
	})

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	got := buf.String()
	// Expected format: GET /books 201 <N>ms 192.168.1.1
	pattern := `^GET /books 201 \d+ms 192\.168\.1\.1\n$`
	if !regexp.MustCompile(pattern).MatchString(got) {
		t.Errorf("log output %q does not match pattern %q", got, pattern)
	}
}

func TestLoggingDefaultStatus(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(nil)
		log.SetFlags(log.LstdFlags)
	})

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/books/123", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	got := buf.String()
	pattern := `^POST /books/123 200 \d+ms 10\.0\.0\.1\n$`
	if !regexp.MustCompile(pattern).MatchString(got) {
		t.Errorf("log output %q does not match pattern %q", got, pattern)
	}
}
