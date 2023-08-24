package films

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var expected = Film{
	FilmID:      1,
	Title:       "title",
	Description: "description",
	ReleaseYear: 2021,
	Rating:      "rating",
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")

	defer db.Close()

	filmMockRows := sqlmock.NewRows([]string{"film_id", "title", "description", "release_year", "rating"})
	filmMockRows.AddRow(expected.FilmID, expected.Title, expected.Description, expected.ReleaseYear, expected.Rating)

	mock.ExpectQuery(regexp.QuoteMeta(SQL_GET_ALL)).
		WillReturnRows(filmMockRows)

	repo := NewPostgresFilmRepository(db)
	films, err := repo.GetAll(context.Background())

	assert.NoError(t, err, "an error '%s' was not expected while getting films", err)
	assert.NotNil(t, films, "films should not be nil")
	assert.Len(t, films, 1, "we expected one film but got %d", len(films))
	assert.Equal(t, 1, films[0].FilmID, "we expected film id to be 1 but got %d", films[0].FilmID)
	assert.Equal(t, expected, films[0], "we expected film to be %v but got %v", expected, films[0])
	assert.NoError(t, mock.ExpectationsWereMet(), "an error '%s' was not expected while getting films", err)
}

func TestGetAllObservesContext(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")

	defer db.Close()

	repo := NewPostgresFilmRepository(db)
	context, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = repo.GetAll(context)
	assert.Error(t, err, "an error '%s' was expected while getting films", err)
	assert.Contains(t, err.Error(), "context canceled", "we expected error to contain 'context canceled' but got %s", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet(), "an error '%s' was not expected while getting films", err)
}

func TestGetByIDObservesContext(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")

	defer db.Close()

	repo := NewPostgresFilmRepository(db)
	context, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = repo.GetByID(context, 1)
	assert.Error(t, err, "an error '%s' was expected while getting film", err)
	assert.Contains(t, err.Error(), "context canceled", "we expected error to contain 'context canceled' but got %s", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet(), "an error '%s' was not expected while getting films", err)
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection")

	defer db.Close()

	filmMockRows := sqlmock.NewRows([]string{"film_id", "title", "description", "release_year", "rating"})
	filmMockRows.AddRow(expected.FilmID, expected.Title, expected.Description, expected.ReleaseYear, expected.Rating)

	mock.ExpectQuery(regexp.QuoteMeta(SQL_BY_ID)).
		WithArgs(expected.FilmID).
		WillReturnRows(filmMockRows)

	repo := NewPostgresFilmRepository(db)
	film, err := repo.GetByID(context.Background(), expected.FilmID)

	assert.NoError(t, err, "an error '%s' was not expected while getting film", err)
	assert.NotNil(t, film, "film should not be nil")
	assert.Equal(t, expected, film, "we expected film to be %v but got %v", expected, film)
	assert.NoError(t, mock.ExpectationsWereMet(), "an error '%s' was not expected while getting films", err)
}
