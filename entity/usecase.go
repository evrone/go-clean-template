package entity

import (
	"context"

	"github.com/evrone/go-service-template/domain"
)

type useCase struct {
	repository domain.EntityRepository
	publisher  domain.EntityPublisher
}

func NewUseCase(r domain.EntityRepository, p domain.EntityPublisher) domain.EntityUseCase {
	return &useCase{r, p}
}

func (u *useCase) Do(ctx context.Context, entity domain.Entity) error {
	entity, err := u.repository.Get(ctx, entity)
	if err != nil {
		return err
	}

	err = u.publisher.Publish(context.Background(), entity)
	if err != nil {
		return err
	}

	return nil
}
