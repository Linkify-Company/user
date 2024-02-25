package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"user/internal/domain"
)

type Person interface {
	AddPerson(ctx context.Context, tx pgx.Tx, person *domain.Person) (int, error)
	GetPersonByUserID(ctx context.Context, userID int) (*domain.Person, error)
	GetAllPersons(ctx context.Context, tx pgx.Tx, reboot bool) ([]*domain.Person, error)

	AddInCash(person *domain.Person)
}

type Repository struct {
	Person
}

type Transaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
}

func NewRepository() *Repository {
	return &Repository{
		Person: NewPersonRepos(),
	}
}
