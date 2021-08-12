package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/evrone/go-service-template/internal/domain"
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

func (s *TranslationService) History() ([]domain.Translation, error) {
	translations, err := s.repo.GetHistory(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "TranslationService - History - s.repo.GetHistory")
	}

	return translations, nil
}

func (s *TranslationService) Translate(translation domain.Translation) (domain.Translation, error) {
	translation, err := s.webAPI.Translate(translation)
	if err != nil {
		return domain.Translation{}, errors.Wrap(err, "TranslationService - Translate - s.webAPI.Translate")
	}

	if err := s.repo.Store(context.Background(), translation); err != nil {
		return domain.Translation{}, errors.Wrap(err, "TranslationService - Translate - s.repo.Store")
	}

	return translation, nil
}
