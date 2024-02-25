package service

import "errors"

var (
	PersonIsAlreadyExist = errors.New("person is already exist")
	PersonNotExist       = errors.New("person is not exist")
)
