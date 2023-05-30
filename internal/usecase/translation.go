package usecase

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/domain/translation/entity"
	"github.com/evrone/go-clean-template/internal/domain/translation/service"
	"github.com/evrone/go-clean-template/internal/usecase/repository"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
)

// TranslationUseCase -.
type TranslationUseCase struct {
	repo   entity.TranslationRepository
	webAPI service.Translator
}

// New -.
func New(r *repository.TranslationRepository, t *webapi.TranslationWebAPI) *TranslationUseCase {
	return NewWithDependencies(
		r,
		t,
	)
}

func NewWithDependencies(r entity.TranslationRepository, w service.Translator) *TranslationUseCase {
	return &TranslationUseCase{
		repo:   r,
		webAPI: w,
	}
}

// History - getting translate history from store.
func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (uc *TranslationUseCase) Translate(_ context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
