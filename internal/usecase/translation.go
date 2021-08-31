package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/evrone/go-clean-template/internal/entity"
)

type TranslationUseCase struct {
	repo   TranslationRepo
	webAPI TranslationWebAPI
}

func New(r TranslationRepo, w TranslationWebAPI) *TranslationUseCase {
	return &TranslationUseCase{
		repo:   r,
		webAPI: w,
	}
}

func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "TranslationUseCase - History - s.repo.GetHistory")
	}

	return translations, nil
}

func (uc *TranslationUseCase) Translate(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, errors.Wrap(err, "TranslationUseCase - Translate - s.webAPI.Translate")
	}

	err = uc.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, errors.Wrap(err, "TranslationUseCase - Translate - s.repo.Store")
	}

	return translation, nil
}
