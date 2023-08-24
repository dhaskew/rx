package films

import (
	"context"
	"sync"
)

type memFilmRepository struct {
	films []Film
	sync.Mutex
}

func NewMemFilmRepository(films []Film) FilmRepository {
	return &memFilmRepository{
		films: films,
	}
}

func (r *memFilmRepository) GetAll(context context.Context) ([]Film, error) {
	return r.films, nil
}

func (r *memFilmRepository) GetByID(context context.Context, id int) (Film, error) {
	r.Lock()
	defer r.Unlock()
	for _, film := range r.films {
		if film.FilmID == id {
			return film, nil
		}
	}
	return Film{}, ErrNotFound
}

func (r *memFilmRepository) GetAllByRating(context context.Context, rating string) ([]Film, error) {
	r.Lock()
	defer r.Unlock()
	var films []Film
	for _, film := range r.films {
		if film.Rating == rating {
			films = append(films, film)
		}
	}
	return films, nil
}

func (r *memFilmRepository) GetAllByCategory(context context.Context, category string) ([]Film, error) {
	r.Lock()
	defer r.Unlock()
	var films []Film
	for _, film := range r.films {
		if film.Category == category {
			films = append(films, film)
		}
	}
	return films, nil
}
