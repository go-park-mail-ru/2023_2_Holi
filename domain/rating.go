package domain

type Rate struct {
	UserID  int `json:"-"`
	VideoID int `json:"videoId"`
	Rate    int `json:"rate"`
}

type RatingRepository interface {
	Insert(rate Rate) (float64, error)
	Delete(rate Rate) (float64, error)
	Exists(rate Rate) (bool, int, error)
	//SelectRating(videoID int) (float64, error)
}

type RatingUsecase interface {
	Add(rate Rate) (float64, error)
	Remove(rate Rate) (float64, error)
	Rated(rete Rate) (bool, int, error)
}
