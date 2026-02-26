package response

import (
	v1 "github.com/evrone/translation-svc/docs/proto/v1"
	"github.com/evrone/translation-svc/internal/entity"
)

// NewTranslationHistory -.
func NewTranslationHistory(translationHistory entity.TranslationHistory) *v1.GetHistoryResponse {
	history := make([]*v1.TranslationHistory, len(translationHistory.History))

	for i, h := range translationHistory.History {
		history[i] = &v1.TranslationHistory{
			Source:      h.Source,
			Destination: h.Destination,
			Original:    h.Original,
			Translation: h.Translation,
		}
	}

	return &v1.GetHistoryResponse{History: history}
}
