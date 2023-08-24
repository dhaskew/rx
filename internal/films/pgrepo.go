package films

import (
	"context"
	"database/sql"
	"errors"
)

const (
	SQL_BY_ID           = `SELECT film_id, title, description , release_year, rating FROM film WHERE film_id = $1 ORDER BY title ASC`
	SQL_GET_ALL         = `SELECT film_id, title, description , release_year, rating FROM film ORDER BY title ASC`
	SQL_GET_BY_RATING   = `SELECT film_id, title, description , release_year FROM film WHERE rating = $1 ORDER BY title ASC`
	SQL_GET_BY_CATEGORY = `SELECT film_id, title, description , release_year, rating from film where film_id in(select distinct(film_id) from film_category where category_id = (select category_id from category where category.name = $1)) ORDER BY title ASC`
)

var ErrNotFound = errors.New("film not found")

type postgressFilmRepository struct {
	db *sql.DB
}

func NewPostgresFilmRepository(db *sql.DB) FilmRepository {
	return &postgressFilmRepository{
		db: db,
	}
}

func (r *postgressFilmRepository) GetAll(context context.Context) ([]Film, error) {
	rows, err := r.db.QueryContext(context, SQL_GET_ALL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear, &film.Rating)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *postgressFilmRepository) GetByID(context context.Context, id int) (Film, error) {
	var film Film
	err := r.db.QueryRowContext(context, SQL_BY_ID, id).Scan(&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear, &film.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return Film{}, ErrNotFound
		}
		return Film{}, err
	}
	return film, nil
}

func (r *postgressFilmRepository) GetAllByRating(context context.Context, rating string) ([]Film, error) {
	rows, err := r.db.QueryContext(context, SQL_GET_BY_RATING, rating)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear)
		if err != nil {
			return nil, err
		}
		film.Rating = rating
		films = append(films, film)
	}
	return films, nil
}

func (r *postgressFilmRepository) GetAllByCategory(context context.Context, category string) ([]Film, error) {
	rows, err := r.db.QueryContext(context, SQL_GET_BY_CATEGORY, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear, &film.Rating)
		if err != nil {
			return nil, err
		}
		film.Category = category
		films = append(films, film)
	}
	return films, nil
}
