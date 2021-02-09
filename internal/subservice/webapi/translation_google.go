package webapi

import (
	translator "github.com/Conight/go-googletrans"
	"github.com/evrone/go-service-template/internal/domain"
)

type GoogleTranslator struct{}

func NewGoogleTranslator() *GoogleTranslator {
	return &GoogleTranslator{}
}

func (y *GoogleTranslator) Translate(entity domain.Translation) (domain.Translation, error) {
	c := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.as"},
	}
	t := translator.New(c)
	result, err := t.Translate(entity.Original, entity.Source, entity.Destination)
	if err != nil {
		return domain.Translation{}, err
	}

	entity.Translation = result.Text

	return entity, nil
}
