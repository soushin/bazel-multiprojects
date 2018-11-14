package usecase

import (
		"github.com/google/go-cloud/wire"
)

var GreetUsecaseSet = wire.NewSet(ProvideUseCase)

type GreetUsecase struct {
	Greet string
}

func ProvideUseCase(greet string) GreetUsecase {
	return GreetUsecase{
		Greet: greet,
	}
}
