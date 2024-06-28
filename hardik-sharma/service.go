package main

import (
	"errors"
	"log"
)

var ErrInvalidId = errors.New("invalid id")
var ErrInvalidContactNo = errors.New("invalid contact number")

type CustomerService interface {
	addCustomer(Customer) error
	updateCustomer(Customer) error
	getAllCustomer() ([]Customer, error)
	getCustomerById(id string) (Customer, error)
	deleteCustomer(id string) error
	subscribe(s Subscriber)
	unSubscribe(s Subscriber)
}

type Service struct {
	customerRepo   Repo
	subscriberList []Subscriber
}

func NewService(repo Repo) *Service {
	return &Service{customerRepo: repo}
}

func (s *Service) subscribe(subs Subscriber) {
	s.subscriberList = append(s.subscriberList, subs)
}

func (s *Service) unSubscribe(subs Subscriber) {
	for i, subscriber := range s.subscriberList {

		if subscriber.getSubscriberId() == subs.getSubscriberId() {
			s.subscriberList = append(s.subscriberList[:i], s.subscriberList[i+1:]...)
		}
	}
}

func (s *Service) notify() {
	customers, err := s.getAllCustomer()
	if err != nil {
		log.Printf("cant fetch data :%q", err)
		return
	}

	for _, subscriber := range s.subscriberList {
		subscriber.update(customers)
	}
}

func calculateDigits(num int) int {
	count := 0
	for i := 0; num != 0; i++ {
		num = num / 10
		count++
	}
	return count
}

func validateId(id string) error {
	if len(id) != 2 {
		return ErrInvalidId
	}
	return nil
}

func validateContactNo(contactNo int) error {
	if calculateDigits(contactNo) != 10 {
		return ErrInvalidContactNo
	}
	return nil
}

func validateCustomer(customer Customer) error {
	if err := validateId(customer.Id); err != nil {
		return err
	}

	if err := validateContactNo(customer.CustomerDetails.ContactNo); err != nil {
		return err
	}

	return nil
}

func (s *Service) addCustomer(customer Customer) error {
	if err := validateCustomer(customer); err != nil {
		return err
	}

	if err := s.customerRepo.create(customer); err != nil {
		return err
	}

	s.notify()
	return nil
}

func (s *Service) updateCustomer(customer Customer) error {
	if err := validateCustomer(customer); err != nil {
		return err
	}

	if err := s.customerRepo.update(customer.Id, customer); err != nil {
		return err
	}

	s.notify()
	return nil
}

func (s *Service) getAllCustomer() ([]Customer, error) {
	return s.customerRepo.getAll()
}

func (s *Service) getCustomerById(id string) (Customer, error) {
	if err := validateId(id); err != nil {
		return Customer{}, err
	}

	return s.customerRepo.getById(id)
}

func (s *Service) deleteCustomer(id string) error {
	if err := validateId(id); err != nil {
		return err
	}

	if err := s.customerRepo.delete(id); err != nil {
		return err
	}

	s.notify()
	return nil
}
