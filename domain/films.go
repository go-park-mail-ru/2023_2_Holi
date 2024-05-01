package domain

type Cast struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Brithday string `json:"brithday"`
	Place    string `json:"place"`
	Carier   string `json:"carier"`
	ImgPath  string `json:"imgPath"`
}

type FilmsRepository interface {
	GetFilmsByGenre(genre int) ([]Video, error)
	GetFilmData(id int) (Video, error)
	GetFilmCast(filmId int) ([]Cast, error)
	GetCastPage(id int) ([]Video, error)
	GetCastName(id int) (Cast, error)
	GetTopRate() (Video, error)
}

type FilmsUsecase interface {
	GetFilmsByGenre(genre int) ([]Video, error)
	GetFilmData(id int) (Video, []Cast, error)
	GetCastPage(id int) ([]Video, Cast, error)
	GetTopRate() (Video, error)
	GetRecommendations(userID int) ([]Video, error)
}

type RecomRepository interface {
	GetRecommendations(userId int) ([]int, error)
}
