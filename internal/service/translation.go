package service

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("TranslationService - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

func (s *TranslationService) Translate(translation domain.Translation) (domain.Translation, error) {
	translation, err := s.webAPI.Translate(translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("TranslationService - Translate - s.webAPI.Translate: %w", err)
	}

	err = s.repo.Store(context.Background(), translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("TranslationService - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
