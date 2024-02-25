package service

import (
	"context"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"time"
	"user/internal/repository"
)

type RebootService struct {
	log         logger.Logger
	transaction repository.Transaction
	personRepos repository.Person
}

func NewRebootService(
	log logger.Logger,
	transaction repository.Transaction,
	personRepos repository.Person,
) Reboot {
	return &RebootService{
		log:         log,
		transaction: transaction,
		personRepos: personRepos,
	}
}

func (m *RebootService) RebootSystem(ctx context.Context) errify.IError {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "Reboot/Begin")
	}

	_, err = m.personRepos.GetAllPersons(ctx, tx, true)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "Reboot/GetAllPersons")
	}

	m.log.Infof("%s / Reboot system successfully", time.Now().Format(time.DateTime))
	return nil
}
