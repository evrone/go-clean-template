package webapi

import (
	translator "github.com/Conight/go-googletrans"
	"github.com/pkg/errors"

	"github.com/evrone/go-clean-template/internal/entity"
)

type TranslationWebAPI struct {
	conf translator.Config
}

func NewTranslationWebAPI() *TranslationWebAPI {
	conf := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &TranslationWebAPI{conf}
}

func (t *TranslationWebAPI) Translate(translation entity.Translation) (entity.Translation, error) {
	trans := translator.New(t.conf)

	result, err := trans.Translate(translation.Original, translation.Source, translation.Destination)
	if err != nil {
		return entity.Translation{}, errors.Wrap(err, "TranslationWebAPI - Translate - trans.Translate")
	}

	translation.Translation = result.Text

	return translation, nil
}
