package service

import "github.com/evrone/go-service-template/internal/domain"

type Translation interface {
	DoTranslate(entity domain.Translation) (domain.Translation, error)
	History() ([]domain.Translation, error)
}
