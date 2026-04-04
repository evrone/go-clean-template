package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase/translation"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

func newTranslationUseCase(t *testing.T) (*translation.UseCase, *MockTranslationRepo, *MockTranslationWebAPI) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := NewMockTranslationRepo(ctrl)
	webAPI := NewMockTranslationWebAPI(ctrl)

	useCase := translation.New(repo, webAPI)

	return useCase, repo, webAPI
}

func TestHistory(t *testing.T) {
	t.Parallel()

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		uc, repo, _ := newTranslationUseCase(t)
		repo.EXPECT().GetHistory(context.Background(), "").Return(nil, nil)

		res, err := uc.History(context.Background(), "")

		require.Equal(t, entity.TranslationHistory{}, res)
		require.NoError(t, err)
	})

	t.Run("result with error", func(t *testing.T) {
		t.Parallel()

		uc, repo, _ := newTranslationUseCase(t)
		repo.EXPECT().GetHistory(context.Background(), "").Return(nil, errInternalServErr)

		res, err := uc.History(context.Background(), "")

		require.Equal(t, entity.TranslationHistory{}, res)
		require.ErrorIs(t, err, errInternalServErr)
	})
}

func TestTranslate(t *testing.T) {
	t.Parallel()

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		uc, repo, webAPI := newTranslationUseCase(t)
		webAPI.EXPECT().Translate(context.Background(), entity.Translation{}).Return(entity.Translation{}, nil)
		repo.EXPECT().Store(context.Background(), "", entity.Translation{}).Return(nil)

		res, err := uc.Translate(context.Background(), "", entity.Translation{})

		require.EqualValues(t, entity.Translation{}, res)
		require.NoError(t, err)
	})

	t.Run("web API error", func(t *testing.T) {
		t.Parallel()

		uc, _, webAPI := newTranslationUseCase(t)
		webAPI.EXPECT().Translate(context.Background(), entity.Translation{}).Return(entity.Translation{}, errInternalServErr)

		res, err := uc.Translate(context.Background(), "", entity.Translation{})

		require.EqualValues(t, entity.Translation{}, res)
		require.ErrorIs(t, err, errInternalServErr)
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		uc, repo, webAPI := newTranslationUseCase(t)
		webAPI.EXPECT().Translate(context.Background(), entity.Translation{}).Return(entity.Translation{}, nil)
		repo.EXPECT().Store(context.Background(), "", entity.Translation{}).Return(errInternalServErr)

		res, err := uc.Translate(context.Background(), "", entity.Translation{})

		require.EqualValues(t, entity.Translation{}, res)
		require.ErrorIs(t, err, errInternalServErr)
	})
}
