package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/evrone/go-clean-template/internal/domain"
)

type TranslationService struct {
	repo   TranslationRepo
	webAPI TranslationWebAPI
}

func NewTranslationService(r TranslationRepo, w TranslationWebAPI) *TranslationService {
	return &TranslationService{
		repo:   r,
		webAPI: w,
	}
}

func (s *TranslationService) History(ctx context.Context) ([]domain.Translation, error) {
	translations, err := s.repo.GetHistory(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "TranslationService - History - s.repo.GetHistory")
	}

	return translations, nil
}

func (s *TranslationService) Translate(ctx context.Context, t domain.Translation) (domain.Translation, error) {
	translation, err := s.webAPI.Translate(t)
	if err != nil {
		return domain.Translation{}, errors.Wrap(err, "TranslationService - Translate - s.webAPI.Translate")
	}

	err = s.repo.Store(context.Background(), translation)
	if err != nil {
		return domain.Translation{}, errors.Wrap(err, "TranslationService - Translate - s.repo.Store")
	}

	return translation, nil
}
