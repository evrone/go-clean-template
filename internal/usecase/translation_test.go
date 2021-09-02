package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func translation(t *testing.T) (*usecase.TranslationUseCase, *MockTranslationRepo, *MockTranslationWebAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockTranslationRepo(mockCtl)
	webAPI := NewMockTranslationWebAPI(mockCtl)

	translation := usecase.New(repo, webAPI)

	return translation, repo, webAPI
}

func TestHistory(t *testing.T) {
	t.Parallel()

	translation, repo, _ := translation(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().GetHistory(context.Background()).Return(nil, nil)
			},
			res: []entity.Translation(nil),
			err: nil,
		},
		{
			name: "result with error",
			mock: func() {
				repo.EXPECT().GetHistory(context.Background()).Return(nil, errInternalServErr)
			},
			res: []entity.Translation(nil),
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := translation.History(context.Background())

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestTranslate(t *testing.T) {
	t.Parallel()

	translation, repo, webAPI := translation(t)

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

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			res, err := translation.Translate(context.Background(), entity.Translation{})

			require.EqualValues(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
