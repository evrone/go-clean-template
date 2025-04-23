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

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func translationUseCase(t *testing.T) (*translation.UseCase, *MockTranslationRepo, *MockTranslationWebAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockTranslationRepo(mockCtl)
	webAPI := NewMockTranslationWebAPI(mockCtl)

	useCase := translation.New(repo, webAPI)

	return useCase, repo, webAPI
}

func TestHistory(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	translationUseCase, repo, _ := translationUseCase(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().GetHistory(context.Background()).Return(nil, nil)
			},
			res: entity.TranslationHistory{},
			err: nil,
		},
		{
			name: "result with error",
			mock: func() {
				repo.EXPECT().GetHistory(context.Background()).Return(nil, errInternalServErr)
			},
			res: entity.TranslationHistory{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := translationUseCase.History(context.Background())

			require.Equal(t, res, localTc.res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTranslate(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	translationUseCase, repo, webAPI := translationUseCase(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, nil)
				repo.EXPECT().Store(context.Background(), entity.Translation{}).Return(nil)
			},
			res: entity.Translation{},
			err: nil,
		},
		{
			name: "web API error",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
		{
			name: "repo error",
			mock: func() {
				webAPI.EXPECT().Translate(entity.Translation{}).Return(entity.Translation{}, nil)
				repo.EXPECT().Store(context.Background(), entity.Translation{}).Return(errInternalServErr)
			},
			res: entity.Translation{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := translationUseCase.Translate(context.Background(), entity.Translation{})

			require.EqualValues(t, res, localTc.res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
