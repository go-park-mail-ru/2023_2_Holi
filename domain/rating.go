package domain

type Rate struct {
	UserID  int `json:"-"`
	VideoID int `json:"videoId"`
	Rate    int `json:"rate"`
}

type RatingRepository interface {
	Insert(rate Rate) error
	Delete(rate Rate) error
	Exists(rate Rate) (bool, error)
	SelectRating(videID int) (int, error)
}

type RatingUsecase interface {
	Add(rate Rate) (int, error)
	Remove(rate Rate) (int, error)
	Rated(rete Rate) (bool, error)
}
