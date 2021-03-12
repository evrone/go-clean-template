package service

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/internal/domain"
	"github.com/evrone/go-service-template/internal/repository"
	"github.com/evrone/go-service-template/internal/webapi"
)

type TranslationService struct {
	repository repository.Translation
	webAPI     webapi.Translation
}

func NewTranslationService(repo repository.Translation, webAPI webapi.Translation) *TranslationService {
	return &TranslationService{
		repository: repo,
		webAPI:     webAPI,
	}
}

func (s *TranslationService) History() ([]domain.Translation, error) {
	translations, err := s.repository.GetHistory(context.Background())
	if err != nil {
		return nil, fmt.Errorf("TranslationService - History - s.repository.GetHistory: %w", err)
	}

	return translations, nil
}

func (s *TranslationService) Translate(translation domain.Translation) (domain.Translation, error) {
	translation, err := s.webAPI.Translate(translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("TranslationService - Translate - s.webAPI.Translate: %w", err)
	}

	err = s.repository.Store(context.Background(), translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("TranslationService - Translate - s.repository.Store: %w", err)
	}

	return translation, nil
}
