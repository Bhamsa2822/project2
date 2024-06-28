package main

import "errors"

var ErrConflict = errors.New("customer already exists")
var ErrNotFound = errors.New("customer not found")

type Repo interface {
	create(c Customer) error
	getAll() ([]Customer, error)
	getById(id string) (Customer, error)
	update(id string, updateCustomer Customer) error
	delete(id string) error
}

type InMemoryRepo struct {
	customers []Customer
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{customers: []Customer{}}
}

func (m *InMemoryRepo) create(newCustomer Customer) error {
	for _, existingCustomer := range m.customers {
		if newCustomer.Id == existingCustomer.Id {
			return ErrConflict
		}
	}

	m.customers = append(m.customers, newCustomer)
	return nil
}

func (m *InMemoryRepo) getAll() ([]Customer, error) {
	return m.customers, nil
}

func (m *InMemoryRepo) getById(id string) (Customer, error) {
	for _, existingCustomer := range m.customers {
		if id == existingCustomer.Id {
			return existingCustomer, nil
		}
	}

	return Customer{}, ErrNotFound
}

func (m *InMemoryRepo) update(id string, updateCustomer Customer) error {
	for i, existingCustomer := range m.customers {
		if existingCustomer.Id == id {
			m.customers[i] = updateCustomer
			return nil
		}
	}
	return ErrNotFound
}

func (m *InMemoryRepo) delete(id string) error {
	for i, existingCustomer := range m.customers {
		if existingCustomer.Id == id {
			m.customers = append(m.customers[:i], m.customers[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
