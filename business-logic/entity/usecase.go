package entity

import (
	"context"
	"log"

	"github.com/evrone/go-service-template/business-logic/domain"
)

type useCase struct {
	repository domain.EntityRepository
	translator domain.EntityTranslator
}

func NewUseCase(repo domain.EntityRepository, api domain.EntityTranslator) domain.EntityUseCase {
	return &useCase{repo, api}
}

func (u *useCase) History() ([]domain.Entity, error) {
	entities, err := u.repository.GetHistory(context.Background())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return entities, nil
}

func (u *useCase) DoTranslate(entity domain.Entity) (domain.Entity, error) {
	entity, err := u.translator.Translate(entity)
	if err != nil {
		log.Print(err)
		return domain.Entity{}, err
	}

	err = u.repository.Store(context.Background(), entity)
	if err != nil {
		log.Print(err)
		return domain.Entity{}, err
	}

	//keys, err := u.repository.GetHistory(context.Background())
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//
	//for _, key := range keys {
	//	fmt.Println(key.Original, key.Translation)
	//}

	return entity, nil
}
