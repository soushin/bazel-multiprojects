package usecase

import (
	"github.com/google/go-cloud/wire"
)

var GreetUsecaseSet = wire.NewSet(ProvideUseCase)

type GreetUsecase struct {
	Msg string
}

func ProvideUseCase(msg string) GreetUsecase {
	return GreetUsecase{
		Msg: msg,
	}
}
