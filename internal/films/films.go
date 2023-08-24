package films

type Film struct {
	FilmID      int    `json:"film_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ReleaseYear int    `json:"release_year,omitempty"`
	Rating      string `json:"rating,omitempty"`
	Category    string `json:"category,omitempty"`
}
