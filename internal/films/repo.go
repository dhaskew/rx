package films

import (
	"context"
)

type FilmRepository interface {
	GetAll(context.Context) ([]Film, error)
	GetByID(context.Context, int) (Film, error)
	GetAllByRating(context.Context, string) ([]Film, error)
	GetAllByCategory(context.Context, string) ([]Film, error)
}
