package entity

import (
	"context"

	"github.com/evrone/go-service-template/internal/business-logic/domain"
)

type useCase struct {
	repo       domain.EntityRepository
	translator domain.EntityTranslator
}

func NewUseCase(
	repository domain.EntityRepository,
	translateAPI domain.EntityTranslator,
) domain.EntityUseCase {
	return &useCase{
		repo:       repository,
		translator: translateAPI,
	}
}

func (u *useCase) History() ([]domain.Entity, error) {
	entities, err := u.repo.GetHistory(context.Background())
	if err != nil {
		domain.Logger.Error(err, "History - repo.GetHistory")
		return nil, err
	}

	return entities, nil
}

func (u *useCase) DoTranslate(entity domain.Entity) (domain.Entity, error) {
	entity, err := u.translator.Translate(entity)
	if err != nil {
		domain.Logger.Error(err, "DoTranslate - translator.Translate")
		return domain.Entity{}, err
	}

	err = u.repo.Store(context.Background(), entity)
	if err != nil {
		domain.Logger.Error(err, "DoTranslate - repo.Store")
		return domain.Entity{}, err
	}

	return entity, nil
}
