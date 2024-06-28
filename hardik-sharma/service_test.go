package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSubscriber struct {
	id           string
	customerList []Customer
}

func newMockSubscriber(id string) *mockSubscriber {
	return &mockSubscriber{
		id:           id,
		customerList: []Customer{},
	}
}

func (m *mockSubscriber) update(customers []Customer) {
	m.customerList = append(m.customerList, customers...)
	return
}

func (m *mockSubscriber) getSubscriberId() string {
	return m.id
}

func TestService_addCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	type args struct {
		newCustomer Customer
	}

	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantCustomers         []Customer
		wantErr               error
		isSubscriber          bool
		wantNotifiedCustomers []Customer
	}{
		{
			name: "invalid id without subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				newCustomer: Customer{
					Id: "hsv",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrInvalidId,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "invalid contact number without subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				newCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 999999,
					},
				},
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrInvalidContactNo,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "valid customer without subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				newCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantErr:               nil,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "adding existing customer without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				newCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr:               ErrConflict,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "valid customer with subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				newCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantErr:      nil,
			isSubscriber: true,
			wantNotifiedCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
		},
		{
			name: "adding existing customer with subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				newCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr:               ErrConflict,
			isSubscriber:          true,
			wantNotifiedCustomers: []Customer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}
			service := &Service{customerRepo: repo}
			subscriber1 := newMockSubscriber("1")

			if tt.isSubscriber {
				service.subscribe(subscriber1)
			}

			gotErr := service.addCustomer(tt.args.newCustomer)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expected error to be same")

			assert.Equal(t, tt.wantCustomers, repo.customers, "expect customers to be matched")

			assert.Equal(t, tt.wantNotifiedCustomers, subscriber1.customerList, "expect customer list to be matched")
		})
	}
}

func TestService_updateCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	type args struct {
		updatedCustomer Customer
	}

	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantCustomers         []Customer
		wantErr               error
		isSubscriber          bool
		wantNotifiedCustomers []Customer
	}{
		{
			name: "invalid id without subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hsv",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrInvalidId,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "invalid contact number without subscriber",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 999999,
					},
				},
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrInvalidContactNo,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "valid updation without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "rj",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "param",
							Address:   "rj",
							ContactNo: 9999999999,
						},
					},
				},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "vs",
					CustomerDetails: CustomerDetails{
						Name:      "varshil",
						Address:   "rj",
						ContactNo: 9999999999,
					},
				},
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantErr:               nil,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "updating non existing customer without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hm",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "rj",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr:               ErrNotFound,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "valid updation with subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "rj",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "param",
							Address:   "rj",
							ContactNo: 9999999999,
						},
					},
				},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "vs",
					CustomerDetails: CustomerDetails{
						Name:      "varshil",
						Address:   "rj",
						ContactNo: 9999999999,
					},
				},
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantErr:      nil,
			isSubscriber: true,
			wantNotifiedCustomers: []Customer{
				{
					Id: "vs",
					CustomerDetails: CustomerDetails{
						Name:      "varshil",
						Address:   "rj",
						ContactNo: 9999999999,
					},
				},
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
		},
		{
			name: "updating non existing customer with subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				updatedCustomer: Customer{
					Id: "hm",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "rj",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr:               ErrNotFound,
			isSubscriber:          true,
			wantNotifiedCustomers: []Customer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}
			service := NewService(repo)
			subscriber1 := newMockSubscriber("1")

			if tt.isSubscriber {
				service.subscribe(subscriber1)
			}

			gotErr := service.updateCustomer(tt.args.updatedCustomer)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expected errors to be same")

			assert.Equal(t, tt.wantCustomers, repo.customers, "expected customer list to be same")

			assert.Equal(t, tt.wantNotifiedCustomers, subscriber1.customerList, "expected customer list to be same")
		})
	}
}

func TestService_getAll(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	tests := []struct {
		name          string
		fields        fields
		wantCustomers []Customer
		wantErr       error
	}{
		{
			name: "getting all customers",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "rj",
							ContactNo: 9999999999,
						},
					},
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "rj",
							ContactNo: 8888888888,
						},
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "rj",
						ContactNo: 9999999999,
					},
				},
				{
					Id: "vs",
					CustomerDetails: CustomerDetails{
						Name:      "varshil",
						Address:   "rj",
						ContactNo: 8888888888,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "test for empty list",
			fields: fields{
				customers: []Customer{},
			},
			wantCustomers: []Customer{},
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)

			gotCustomers, gotErr := service.getAllCustomer()

			assert.ErrorIs(t, tt.wantErr, gotErr, "expected error to be matched")

			assert.Equal(t, tt.wantCustomers, gotCustomers, "expected customer list to be matched")
		})
	}
}

func TestService_getCustomerById(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	type args struct {
		id string
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantCustomer Customer
		wantErr      error
	}{
		{
			name: "invalid id",
			fields: fields{
				customers: []Customer{
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "udaipur",
							ContactNo: 9999999959,
						},
					},
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 8619185565,
						},
					},
				},
			},
			args: args{
				id: "hsvv",
			},
			wantCustomer: Customer{},
			wantErr:      ErrInvalidId,
		},
		{
			name: "valid id and existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "udaipur",
							ContactNo: 9999999959,
						},
					},
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 8619185565,
						},
					},
				},
			},
			args: args{
				id: "hs",
			},
			wantCustomer: Customer{
				Id: "hs",
				CustomerDetails: CustomerDetails{
					Name:      "hardik",
					Address:   "udaipur",
					ContactNo: 8619185565,
				},
			},
			wantErr: nil,
		},
		{
			name: "non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "vs",
						CustomerDetails: CustomerDetails{
							Name:      "varshil",
							Address:   "udaipur",
							ContactNo: 9999999959,
						},
					},
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 8619185565,
						},
					},
				},
			},
			args: args{
				id: "vv",
			},
			wantCustomer: Customer{},
			wantErr:      ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)

			gotCustomer, gotErr := service.getCustomerById(tt.args.id)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expected error to be same")

			assert.Equal(t, tt.wantCustomer, gotCustomer, "expected customer to be same")
		})
	}
}

func TestService_deleteCustomer(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	type args struct {
		id string
	}

	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantCustomers         []Customer
		wantErr               error
		isSubscriber          bool
		wantNotifiedCustomers []Customer
	}{
		{
			name: "invalid id without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "hf",
						CustomerDetails: CustomerDetails{
							Name:      "hk",
							Address:   "udr",
							ContactNo: 8888888888,
						},
					},
				},
			},
			args: args{
				id: "hsss",
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 7777777777,
					},
				},
				{
					Id: "hf",
					CustomerDetails: CustomerDetails{
						Name:      "hk",
						Address:   "udr",
						ContactNo: 8888888888,
					},
				},
			},
			wantErr:               ErrInvalidId,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "deleting non existing customer without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "hf",
						CustomerDetails: CustomerDetails{
							Name:      "hk",
							Address:   "udr",
							ContactNo: 8888888888,
						},
					},
				},
			},
			args: args{
				id: "hm",
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 7777777777,
					},
				},
				{
					Id: "hf",
					CustomerDetails: CustomerDetails{
						Name:      "hk",
						Address:   "udr",
						ContactNo: 8888888888,
					},
				},
			},
			wantErr:               ErrNotFound,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "deleting customer without subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "hf",
						CustomerDetails: CustomerDetails{
							Name:      "hk",
							Address:   "udr",
							ContactNo: 8888888888,
						},
					},
				},
			},
			args: args{
				id: "hs",
			},
			wantCustomers: []Customer{
				{
					Id: "hf",
					CustomerDetails: CustomerDetails{
						Name:      "hk",
						Address:   "udr",
						ContactNo: 8888888888,
					},
				},
			},
			wantErr:               nil,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name:   "deleting from empty database without subscriber",
			fields: fields{customers: []Customer{}},
			args: args{
				id: "hs",
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrNotFound,
			isSubscriber:          false,
			wantNotifiedCustomers: []Customer{},
		},
		{
			name: "deleting customer with subscriber",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 7777777777,
						},
					},
					{
						Id: "hf",
						CustomerDetails: CustomerDetails{
							Name:      "hk",
							Address:   "udr",
							ContactNo: 8888888888,
						},
					},
				},
			},
			args: args{
				id: "hs",
			},
			wantCustomers: []Customer{
				{
					Id: "hf",
					CustomerDetails: CustomerDetails{
						Name:      "hk",
						Address:   "udr",
						ContactNo: 8888888888,
					},
				},
			},
			wantErr:      nil,
			isSubscriber: true,
			wantNotifiedCustomers: []Customer{
				{
					Id: "hf",
					CustomerDetails: CustomerDetails{
						Name:      "hk",
						Address:   "udr",
						ContactNo: 8888888888,
					},
				},
			},
		},
		{
			name:   "deleting from empty database with subscriber",
			fields: fields{customers: []Customer{}},
			args: args{
				id: "hs",
			},
			wantCustomers:         []Customer{},
			wantErr:               ErrNotFound,
			isSubscriber:          true,
			wantNotifiedCustomers: []Customer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{tt.fields.customers}
			service := NewService(repo)
			subscriber1 := newMockSubscriber("1")

			if tt.isSubscriber {
				service.subscribe(subscriber1)
			}

			gotErr := service.deleteCustomer(tt.args.id)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expected error to be same")

			assert.Equal(t, tt.wantCustomers, repo.customers, "expected customer list to be same")

			assert.Equal(t, tt.wantNotifiedCustomers, subscriber1.customerList, "expected customer list to be same")
		})
	}
}

func TestService_subscribe(t *testing.T) {
	subscriber1 := newMockSubscriber("1")
	subscriber2 := newMockSubscriber("2")
	subscriber3 := newMockSubscriber("3")
	subscriber4 := newMockSubscriber("4")

	type fields struct {
		subscribers []Subscriber
	}

	tests := []struct {
		name        string
		fields      fields
		subscribers []Subscriber
		wantLen     int
	}{
		{
			name: "adding 3 subscriber to empty subscriber list",
			fields: fields{
				subscribers: []Subscriber{},
			},
			subscribers: []Subscriber{
				subscriber1,
				subscriber2,
				subscriber3,
			},
			wantLen: 3,
		},
		{
			name: "adding 2 subscriber to subscriber list with 2 existing subscriber",
			fields: fields{
				subscribers: []Subscriber{
					subscriber1,
					subscriber2,
				},
			},
			subscribers: []Subscriber{
				subscriber2,
				subscriber4,
			},
			wantLen: 4,
		},
		{
			name: "not adding subscribers",
			fields: fields{
				subscribers: []Subscriber{
					subscriber1,
					subscriber2,
				},
			},
			subscribers: []Subscriber{},
			wantLen:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepo()
			service := &Service{
				customerRepo:   repo,
				subscriberList: tt.fields.subscribers}

			for _, subscriber := range tt.subscribers {
				service.subscribe(subscriber)
			}

			assert.Equal(t, tt.wantLen, len(service.subscriberList), "expected length to be same")
		})
	}
}

func TestService_unSubscribe(t *testing.T) {
	subscriber2 := newMockSubscriber("1")
	subscriber1 := newMockSubscriber("2")
	subscriber3 := newMockSubscriber("3")

	type fields struct {
		subscribers []Subscriber
	}

	tests := []struct {
		name         string
		fields       fields
		unSubscriber Subscriber
		wantLen      int
	}{
		{
			name: "unsubscribing 1 subscriber",
			fields: fields{
				subscribers: []Subscriber{
					subscriber1,
					subscriber2,
					subscriber3,
				},
			},
			unSubscriber: subscriber1,
			wantLen:      2,
		},
		{
			name: "unsubscribing unknown subscriber",
			fields: fields{
				subscribers: []Subscriber{
					subscriber1,
					subscriber2,
				},
			},
			unSubscriber: subscriber3,
			wantLen:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepo()
			service := &Service{
				customerRepo:   repo,
				subscriberList: tt.fields.subscribers,
			}

			service.unSubscribe(tt.unSubscriber)

			assert.Equal(t, tt.wantLen, len(service.subscriberList), "expected length to be same")
		})
	}
}

func Test_calculateDigits(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
	}{
		{
			name: "10 digit number",
			args: args{
				num: 9649127550,
			},
			wantCount: 10,
		},
		{
			name: "10 digit number having 0",
			args: args{
				num: 9602201100,
			},
			wantCount: 10,
		},
		{
			name: "having all digits 0",
			args: args{
				num: 000000000,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount := calculateDigits(tt.args.num)

			assert.Equal(t, tt.wantCount, gotCount, "expected count to be equal")
		})
	}
}
