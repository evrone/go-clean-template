package webapi_test

import (
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	translation := webapi.New()

	require.IsType(t, translation, &webapi.TranslationWebAPI{})
}

func TestTranslateErr(t *testing.T) {
	t.Parallel()

	translation := webapi.New()

	_, err := translation.Translate(entity.Translation{})

	require.Error(t, err)
}
