package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"user/internal/domain"
	"user/internal/repository/cash"
	"user/internal/repository/postgres"
)

type PersonRepos struct {
	memoryPerson *cash.PersonCash
}

func NewPersonRepos() Person {
	return &PersonRepos{
		memoryPerson: cash.NewPersonCash(),
	}
}

func (m *PersonRepos) GetAllPersons(ctx context.Context, tx pgx.Tx, reboot bool) ([]*domain.Person, error) {
	if reboot {
		rows, err := tx.Query(ctx, `SELECT id, user_id, email, name, patronymic, surname, sex, age, birthday FROM person`)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("GetAllPersons/Query: %w", err)
		}
		defer rows.Close()

		var persons []*domain.Person
		for rows.Next() {
			persons = append(persons, &domain.Person{})
			err = rows.Scan(persons[len(persons)-1].ScanValues()...)
			if err != nil {
				return nil, fmt.Errorf("GetAllPersons/Scan: %w", err)
			}
			m.memoryPerson.Set(persons[len(persons)-1])
		}
		return persons, nil
	}
	return m.memoryPerson.GetAllPersons(), nil
}

func (m *PersonRepos) AddPerson(ctx context.Context, tx pgx.Tx, person *domain.Person) (int, error) {
	if m.memoryPerson.Get(*person.UserID) != nil {
		return 0, PersonAlreadyExist
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO person 
		    (user_id, email, name, patronymic, surname, sex, age, birthday)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
`, person.UserID, person.Email, person.Name, person.Patronymic, person.Surname, person.Sex, person.Age, person.Birthday)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == postgres.ErrUniqueViolation {
			return 0, PersonAlreadyExist
		}
		return 0, fmt.Errorf("AddPerson/Scan: %w", err)
	}
	return id, nil
}

func (m *PersonRepos) GetPersonByUserID(ctx context.Context, userID int) (*domain.Person, error) {
	person := m.memoryPerson.Get(userID)
	if person == nil {
		return nil, PersonNotExist
	}
	return person, nil
}

func (m *PersonRepos) AddInCash(person *domain.Person) {
	m.memoryPerson.Set(person)
}
