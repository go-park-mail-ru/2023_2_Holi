package domain

type UtilsRepository interface {
	GetIdFromStorage(token string) (int, error)
}

type UtilsUsecase interface {
	GetIdBy(token string) (int, error)
}
