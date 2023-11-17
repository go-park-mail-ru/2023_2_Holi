package domain

type UtilsRepository interface {
	GetIdBy(token string) (int, error)
}

type UtilsUsecase interface {
	GetIdBy(token string) (int, error)
}
