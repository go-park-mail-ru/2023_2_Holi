package domain

type Episode struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PreviewPath string `json:"previewPath"`
	MediaPath   string `json:"mediaPath"`
	Number      int    `json:"number"`
	Season      int    `json:"season"`
}

type SeriesRepository interface {
	GetSeriesByGenre(genre int) ([]Video, error)
	GetSeriesData(id int) (Video, error)
	GetSeriesCast(id int) ([]Cast, error)
	GetSeriesEpisodes(id int) ([]Episode, error)
	GetCastPageSeries(id int) ([]Video, error)
	GetCastNameSeries(id int) (Cast, error)
	GetTopRate() (Video, error)
}

type SeriesUsecase interface {
	GetSeriesByGenre(genre int) ([]Video, error)
	GetSeriesData(id int) (Video, []Cast, []Episode, error)
	GetCastPageSeries(id int) ([]Video, Cast, error)
	GetTopRate() (Video, error)
}
