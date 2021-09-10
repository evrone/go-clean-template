package webapi_test

import (
	"testing"

	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	translation := webapi.New()

	require.IsType(t, translation, &webapi.TranslationWebAPI{})
}
