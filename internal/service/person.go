package service

import (
	"context"
	"errors"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/Linkify-Company/common_utils/pointer"
	"user/internal/domain"
	"user/internal/repository"
)

type PersonService struct {
	log         logger.Logger
	transaction repository.Transaction
	personRepos repository.Person
}

func NewPersonService(
	log logger.Logger,
	transaction repository.Transaction,
	personRepos repository.Person,
) Person {
	return &PersonService{
		log:         log,
		transaction: transaction,
		personRepos: personRepos,
	}
}

func (m *PersonService) AddPerson(ctx context.Context, person *domain.Person) (int, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return 0, errify.NewInternalServerError(err.Error(), "AddPerson/Begin")
	}
	defer m.transaction.Rollback(ctx, tx)

	id, err := m.personRepos.AddPerson(ctx, tx, person)
	if err != nil {
		if errors.Is(err, repository.PersonAlreadyExist) {
			return 0, errify.NewBadRequestError(err.Error(), PersonIsAlreadyExist.Error(), "AddPerson/AddPerson")
		}
		return 0, errify.NewInternalServerError(err.Error(), "AddPerson/AddPerson")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, errify.NewInternalServerError(err.Error(), "Commit transaction failed")
	}
	person.ID = pointer.P(id)
	m.personRepos.AddInCash(person)

	return id, nil
}

func (m *PersonService) GetPersonByUserID(ctx context.Context, userID int) (*domain.Person, errify.IError) {
	person, err := m.personRepos.GetPersonByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.PersonNotExist) {
			return nil, errify.NewBadRequestError(err.Error(), PersonNotExist.Error(), "GetPersonByUserID/GetPersonByUserID")
		}
		return nil, errify.NewInternalServerError(err.Error(), "GetPersonByUserID/GetPersonByUserID")
	}
	return person, nil
}
