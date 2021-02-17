package webapi

import (
	"fmt"

	translator "github.com/Conight/go-googletrans"

	"github.com/evrone/go-service-template/internal/domain"
)

type translationWebAPI struct {
	conf translator.Config
}

func NewTranslationWebAPI() Translation {
	conf := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &translationWebAPI{conf}
}

func (t *translationWebAPI) Translate(translation domain.Translation) (domain.Translation, error) {
	trans := translator.New(t.conf)

	result, err := trans.Translate(translation.Original, translation.Source, translation.Destination)
	if err != nil {
		return domain.Translation{}, fmt.Errorf("translationWebAPI - Translate - trans.Translate: %w", err)
	}

	translation.Translation = result.Text

	return translation, nil
}
