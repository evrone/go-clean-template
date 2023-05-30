package service

import "github.com/evrone/go-clean-template/internal/domain/translation/entity"

type Translator interface {
	Translate(translation entity.Translation) (entity.Translation, error)
}
