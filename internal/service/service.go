package service

// TODO Добавить интерфейсы для слоя репо

type Repository interface {
}

type Service struct {
	Repository Repository
}
