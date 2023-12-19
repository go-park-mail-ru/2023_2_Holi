package domain

import "time"

type SubsUsecase interface {
	CheckSub(userId int) (subUpTo time.Time, status bool, err error)
	Subscribe(userId int) error
	UnSubscribe(userId int) error
}

type SubsRepository interface {
	CheckSub(userId int) (subUpTo time.Time, err error)
	Subscribe(userId int) error
	UnSubscribe(userId int) error
}
