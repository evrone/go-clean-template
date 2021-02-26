package service

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/internal/domain"
	"github.com/evrone/go-service-template/internal/repository"
	"github.com/evrone/go-service-template/internal/webapi"
)

type translationService struct {
	repository repository.Translation
	webAPI     webapi.Translation
}

func NewTranslationService(repo repository.Translation, webAPI webapi.Translation) Translation {
	return &translationService{
		repository: repo,
		webAPI:     webAPI,
	}
}

func (u *translationService) History() ([]domain.Translation, error) {
	translations, err := u.repository.GetHistory(context.Background())
	if err != nil {
		return nil, fmt.Errorf("translationService - History - u.repository.GetHistory: %w", err)
	}

	return translations, nil
}

func (u *translationService) Translate(translation domain.Translation) (domain.Translation, error) {
	translation, err := u.webAPI.Translate(translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("translationService - Translate - u.webAPI.Translate: %w", err)
	}

	err = u.repository.Store(context.Background(), translation)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("translationService - Translate - u.repository.Store: %w", err)
	}

	return translation, nil
}
