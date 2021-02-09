package service

import (
	"context"

	"github.com/evrone/go-service-template/internal/subservice/repository"

	"github.com/evrone/go-service-template/internal/subservice/webapi"

	"github.com/evrone/go-service-template/internal/domain"
)

type useCase struct {
	translationRepo   repository.Translation
	translationWebAPI webapi.Translation
}

func NewUseCase(repo repository.Translation, api webapi.Translation) Translation {
	return &useCase{
		translationRepo:   repo,
		translationWebAPI: api,
	}
}

func (u *useCase) History() ([]domain.Translation, error) {
	entities, err := u.translationRepo.GetHistory(context.Background())
	if err != nil {
		domain.Logger.Error(err, "History - translationRepo.GetHistory",
			domain.Field{Key: "key", Val: "value"},
		)
		return nil, err
	}

	return entities, nil
}

func (u *useCase) DoTranslate(entity domain.Translation) (domain.Translation, error) {
	entity, err := u.translationWebAPI.Translate(entity)
	if err != nil {
		domain.Logger.Error(err, "DoTranslate - translationWebAPI.Translate")
		return domain.Translation{}, err
	}

	err = u.translationRepo.Store(context.Background(), entity)
	if err != nil {
		domain.Logger.Error(err, "DoTranslate - translationRepo.Store")
		return domain.Translation{}, err
	}

	return entity, nil
}
