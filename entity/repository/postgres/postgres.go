package postgres

import (
	"context"
	"fmt"
	"github.com/evrone/go-service-template/domain"
)

type Repository struct {
	Connect string
}

func NewEntityRepository(Connect string) domain.EntityRepository {
	return &Repository{Connect}
}

func (r *Repository) GetByID(ctx context.Context, ID int) (string, error) {
	if ID == 41 {
		return r.Connect, nil
	}
	return "", fmt.Errorf("ошибка")
}
