package domain

import (
	"fmt"
	"regexp"
)

type Date string

func (m *Date) Valid() error {
	if m == nil {
		return fmt.Errorf("date is empty")
	}
	pattern := `^\d{4}\.\d{2}\.\d{2}$`
	valid, err := regexp.MatchString(pattern, string(*m))
	if err != nil || !valid {
		return fmt.Errorf("invalid date format")
	}
	return nil
}
