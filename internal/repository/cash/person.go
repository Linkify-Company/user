package cash

import (
	"sync"
	"user/internal/domain"
)

type PersonCash struct {
	// map[userID]*domain.Person
	persons map[int]*domain.Person
	len     int
	mx      sync.RWMutex
}

func NewPersonCash() *PersonCash {
	return &PersonCash{
		persons: make(map[int]*domain.Person),
		mx:      sync.RWMutex{},
	}
}

func (m *PersonCash) Set(person *domain.Person) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.persons[*person.UserID] == nil {
		m.persons[*person.UserID] = person
		m.len++
	}
}

func (m *PersonCash) Get(userID int) *domain.Person {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.persons[userID]
}

func (m *PersonCash) GetAllPersons() []*domain.Person {
	var persons = make([]*domain.Person, m.len)

	for _, person := range m.persons {
		persons[len(persons)-1] = person
	}
	return persons
}
