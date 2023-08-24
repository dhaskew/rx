package films

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data0 = []Film{}

var data1 = []Film{
	{
		FilmID:      1,
		Title:       "Film 1",
		Description: "Description 1",
		ReleaseYear: 2001,
	},
}

var data2 = []Film{
	{
		FilmID:      1,
		Title:       "Film 1",
		Description: "Description 1",
		ReleaseYear: 2001,
	},
	{
		FilmID:      2,
		Title:       "Film 2",
		Description: "Description 2",
		ReleaseYear: 2002,
	},
}

func TestMemFilmRepositoryGetAll(t *testing.T) {
	type args struct {
		dbstart []Film
	}
	tests := []struct {
		name        string
		args        args
		want        FilmRepository
		expected    []Film
		expectedErr error
	}{
		{
			name: "No Entry",
			args: args{
				dbstart: data0,
			},
			want: &memFilmRepository{
				films: data0,
			},
			expected:    data0,
			expectedErr: nil,
		},
		{
			name: "Single Entry",
			args: args{
				dbstart: data1,
			},
			want: &memFilmRepository{
				films: data1,
			},
			expected:    data1,
			expectedErr: nil,
		},
		{
			name: "Multiple Entries",
			args: args{
				dbstart: data2,
			},
			want: &memFilmRepository{
				films: data2,
			},
			expected:    data2,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemFilmRepository(tt.args.dbstart)
			actual, err := repo.GetAll(context.Background())
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestMemFilmRepositoryGetByID(t *testing.T) {
	type args struct {
		dbstart []Film
	}
	tests := []struct {
		name        string
		args        args
		want        int
		expected    Film
		expectedErr error
	}{
		{
			name: "No Entry",
			args: args{
				dbstart: data0,
			},
			want:        1,
			expected:    Film{},
			expectedErr: ErrNotFound,
		},
		{
			name: "Single Entry",
			args: args{
				dbstart: data1,
			},
			want:     1,
			expected: data1[0],
		},
		{
			name: "Multiple Entries",
			args: args{
				dbstart: data2,
			},
			want:     1,
			expected: data2[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemFilmRepository(tt.args.dbstart)
			actual, err := repo.GetByID(context.Background(), tt.want)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}
