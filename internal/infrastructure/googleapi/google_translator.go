package googleapi

import (
	"fmt"

	translator "github.com/Conight/go-googletrans"

	"github.com/evrone/go-clean-template/internal/domain/translation/entity"
)

// GoogleTranslator -.
type GoogleTranslator struct {
	conf translator.Config
}

// New -.
func New() *GoogleTranslator {
	conf := translator.Config{
		UserAgent:   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"},
		ServiceUrls: []string{"translate.google.com"},
	}

	return &GoogleTranslator{
		conf: conf,
	}
}

// Translate -.
func (t *GoogleTranslator) Translate(translation entity.Translation) (entity.Translation, error) {
	trans := translator.New(t.conf)

	result, err := trans.Translate(translation.Original, translation.Source, translation.Destination)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("GoogleTranslator - Translate - trans.Translate: %w", err)
	}

	translation.Translation = result.Text

	return translation, nil
}
