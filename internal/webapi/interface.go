package webapi

import "github.com/evrone/go-service-template/internal/domain"

type Translation interface {
	Translate(entity domain.Translation) (domain.Translation, error)
}
