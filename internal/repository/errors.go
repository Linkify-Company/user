package repository

import "errors"

var (
	PersonAlreadyExist = errors.New("person already exists")
	PersonNotExist     = errors.New("person not exists")
)
