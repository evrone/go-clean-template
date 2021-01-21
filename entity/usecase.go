package entity

import (
	"context"
	"github.com/evrone/go-service-template/domain"
)

type UseCase struct {
	repository domain.EntityRepository
}

func NewUsecase(repository domain.EntityRepository) domain.EntityUsecase {
	return &UseCase{repository}
}

func (u *UseCase) Get(ctx context.Context, ID int) (string, error) {
	msg, err := u.repository.GetByID(ctx, ID)
	return msg, err
}
