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
}

type RatingUsecase interface {
	Add(rate Rate) error
	Remove(rate Rate) error
	Rated(rete Rate) (bool, error)
}
