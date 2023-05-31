package application

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/internal/domain/translation/entity"
	"github.com/evrone/go-clean-template/internal/domain/translation/service"
	"github.com/evrone/go-clean-template/internal/infrastructure/googleapi"
	"github.com/evrone/go-clean-template/internal/infrastructure/repository"
)

// TranslationUseCase -.
type TranslationUseCase struct {
	translationRepository entity.TranslationRepository
	translator            service.Translator
}

// New -.
func New(r *repository.TranslationRepository, t *googleapi.GoogleTranslator) *TranslationUseCase {
	return NewWithDependencies(
		r,
		t,
	)
}

func NewWithDependencies(translationRepository entity.TranslationRepository, translator service.Translator) *TranslationUseCase {
	return &TranslationUseCase{
		translationRepository: translationRepository,
		translator:            translator,
	}
}

// History - getting translate history from store.
func (uc *TranslationUseCase) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := uc.translationRepository.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationUseCase - History - s.translationRepository.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (uc *TranslationUseCase) Translate(_ context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.translator.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.translator.Translate: %w", err)
	}

	err = uc.translationRepository.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.translationRepository.Store: %w", err)
	}

	return translation, nil
}
