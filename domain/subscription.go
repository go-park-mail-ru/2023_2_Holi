package domain

type SubInfo struct {
	ID       int    `json:"id"`
	SubStart string `json:"sub_start"`
	SubUpTo  string `json:"sub_up_to"`
}

type SubsUsecase interface {
	CheckSub(userId int) error
	GetSubInfo(userId int) (SubInfo, error)
	Subscribe(userId int, flag int) error
	UnSubscribe(userId int) error
}

type SubsRepository interface {
	CheckSub(userId int) error
	GetSubInfo(userId int) (SubInfo, error)
	Subscribe(userId int, flag int) error
	UnSubscribe(userId int) error
}
