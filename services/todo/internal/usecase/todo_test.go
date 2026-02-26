package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/evrone/todo-svc/internal/entity"
	"github.com/evrone/todo-svc/internal/usecase/todo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  any
	err  error
}

func todoUseCase(t *testing.T) (*todo.UseCase, *MockTodoRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockTodoRepo(mockCtl)
	useCase := todo.New(repo)

	return useCase, repo
}

func TestTodoCreate(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	uc, repo := todoUseCase(t)

	sampleTodo := entity.Todo{Title: "Test todo", Priority: "high"}

	tests := []test{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Create(context.Background(), gomock.Any()).Return(sampleTodo, nil)
			},
			res: sampleTodo,
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().Create(context.Background(), gomock.Any()).Return(entity.Todo{}, errInternalServErr)
			},
			res: entity.Todo{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := uc.Create(context.Background(), entity.Todo{Title: "Test todo", Priority: "high"})

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTodoGetByID(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	uc, repo := todoUseCase(t)

	sampleTodo := entity.Todo{ID: 1, Title: "Test todo"}

	tests := []test{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().GetByID(context.Background(), 1).Return(sampleTodo, nil)
			},
			res: sampleTodo,
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().GetByID(context.Background(), 1).Return(entity.Todo{}, errInternalServErr)
			},
			res: entity.Todo{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := uc.GetByID(context.Background(), 1)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTodoList(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	uc, repo := todoUseCase(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().List(context.Background()).Return(nil, nil)
			},
			res: ([]entity.Todo)(nil),
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().List(context.Background()).Return(nil, errInternalServErr)
			},
			res: ([]entity.Todo)(nil),
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := uc.List(context.Background())

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTodoUpdate(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	uc, repo := todoUseCase(t)

	updated := entity.Todo{ID: 1, Title: "Updated", Status: "done"}

	tests := []test{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Update(context.Background(), 1, gomock.Any()).Return(updated, nil)
			},
			res: updated,
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().Update(context.Background(), 1, gomock.Any()).Return(entity.Todo{}, errInternalServErr)
			},
			res: entity.Todo{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := uc.Update(context.Background(), 1, entity.Todo{Title: "Updated", Status: "done"})

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestTodoDelete(t *testing.T) { //nolint:tparallel // data races here
	t.Parallel()

	uc, repo := todoUseCase(t)

	tests := []struct {
		name string
		mock func()
		err  error
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Delete(context.Background(), 1).Return(nil)
			},
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().Delete(context.Background(), 1).Return(errInternalServErr)
			},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			err := uc.Delete(context.Background(), 1)

			require.ErrorIs(t, err, localTc.err)
		})
	}
}

