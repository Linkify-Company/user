package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Person struct {
	ID         *int    `json:"id"`
	UserID     *int    `json:"user_id"`
	Email      *string `json:"email" validate:"required,email"`
	Name       *string `json:"name" validate:"required"`
	Patronymic *string `json:"patronymic"`
	Surname    *string `json:"surname,omitempty" validate:"required"`
	Sex        *string `json:"sex,omitempty" validate:"required,oneof=male female"`
	Age        *int    `json:"age,omitempty" validate:"min=1,max=200"`
	Birthday   *Date   `json:"birthday,omitempty"`
}

func (m *Person) Valid() error {
	if m == nil {
		return errors.New("person is empty")
	}
	if m.Birthday != nil {
		err := m.Birthday.Valid()
		if err != nil {
			return err
		}
	}

	vl := validator.New()
	err := vl.Struct(*m)
	if err != nil {
		return err.(validator.ValidationErrors)[0]
	}
	return nil
}

func (m *Person) ScanValues() []any {
	return []any{&m.ID, &m.UserID, &m.Email, &m.Name, &m.Patronymic, &m.Surname, &m.Sex, &m.Age, &m.Birthday}
}
