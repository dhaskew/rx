package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhaskew/rx/internal/films"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestFilmsHandlerNoFilms(t *testing.T) {
	t.Parallel()

	mem := films.NewMemFilmRepository([]films.Film{})
	log := NewLogger()
	srv := NewServer(
		WithFilmRepository(&mem),
		WithLogger(log),
		WithRouterFunc(chi.NewRouter),
		WithPort("8080"),
	)

	req, err := http.NewRequest("GET", "/films", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.filmsHandler())
	handler.ServeHTTP(rr, req)
	expected := `[]`
	actual := rr.Body.String()
	assert.Equal(t, expected, actual)
}

func TestFilmsHandler(t *testing.T) {
	t.Parallel()

	film := films.Film{
		FilmID:      1,
		Title:       "title",
		Description: "description",
		ReleaseYear: 2021,
	}

	memdb := []films.Film{film}

	mem := films.NewMemFilmRepository(memdb)
	log := NewLogger()
	srv := NewServer(
		WithFilmRepository(&mem),
		WithLogger(log),
		WithRouterFunc(chi.NewRouter),
		WithPort("8080"),
	)

	req, err := http.NewRequest("GET", "/v1/films", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.filmsHandler())
	handler.ServeHTTP(rr, req)
	bytes, _ := json.MarshalIndent(memdb, "", "\t")
	expected := string(bytes)
	actual := rr.Body.String()
	assert.Equal(t, expected, actual)
}
