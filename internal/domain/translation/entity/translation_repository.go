package entity

import "context"

type TranslationRepository interface {
	Store(context.Context, Translation) error
	GetHistory(context.Context) ([]Translation, error)
}
