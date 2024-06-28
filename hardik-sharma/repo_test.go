package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestInMemoryRepo_create(t *testing.T) {
	type fields struct {
		customers []Customer
	}
	type args struct {
		newCustomer Customer
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCustomers []Customer
		wantErr       error
	}{
		{
			name: "adding new customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hm",
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
					Id: "hm",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "adding existing customer",
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
			wantErr: ErrConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}

			gotErr := repo.create(tt.args.newCustomer)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expecting error :%q but got error :%q", tt.wantErr, gotErr)
			}

			if !reflect.DeepEqual(repo.customers, tt.wantCustomers) {
				t.Errorf("customer list should match:\nwant: %+v\ngot: %+v", tt.wantCustomers, repo.customers)
			}
		})
	}
}

func TestInMemoryRepo_delete(t *testing.T) {
	type fields struct {
		customers []Customer
	}
	type args struct {
		id string
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCustomers []Customer
		wantErr       error
	}{
		{
			name: "deleting existing customer",
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
					{
						Id: "hm",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "ahmedabad",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				id: "hs",
			},
			wantCustomers: []Customer{
				{
					Id: "hm",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "ahmedabad",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "deleting non existing customer",
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
				id: "hm",
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
			wantErr: ErrNotFound,
		},
		{
			name: "deleting from empty list",
			fields: fields{
				customers: []Customer{},
			},
			args: args{
				id: "hm",
			},
			wantCustomers: []Customer{},
			wantErr:       ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}

			gotErr := repo.delete(tt.args.id)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("got error %q but expecting error :%q", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(repo.customers, tt.wantCustomers) {
				t.Errorf("customer list should be\n got customers %+v\n but expecting %+v", repo.customers, tt.wantCustomers)
			}
		})
	}
}

func TestInMemoryRepo_update(t *testing.T) {
	type fields struct {
		customers []Customer
	}

	type args struct {
		id             string
		updateCustomer Customer
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCustomers []Customer
		wantErr       error
	}{
		{
			name: "updating existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "uadipur",
							ContactNo: 9649127559,
						},
					},
					{
						Id: "hk",
						CustomerDetails: CustomerDetails{
							Name:      "h",
							Address:   "raj",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				id: "hs",
				updateCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "Ahmedabad",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "Ahmedabad",
						ContactNo: 9649127559,
					},
				},
				{
					Id: "hk",
					CustomerDetails: CustomerDetails{
						Name:      "h",
						Address:   "raj",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "updating non existing customer",
			fields: fields{
				customers: []Customer{
					{
						Id: "hk",
						CustomerDetails: CustomerDetails{
							Name:      "h",
							Address:   "raj",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				id: "hm",
				updateCustomer: Customer{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "Ahmedabad",
						ContactNo: 9649127559,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hk",
					CustomerDetails: CustomerDetails{
						Name:      "h",
						Address:   "raj",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}

			gotErr := repo.update(tt.args.id, tt.args.updateCustomer)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expecting error %q but got error :%q", tt.wantErr, gotErr)
			}

			if !reflect.DeepEqual(repo.customers, tt.wantCustomers) {
				t.Errorf("Customer list should match:\nwant: %+v\ngot: %+v", tt.wantCustomers, repo.customers)
			}
		})
	}
}

func TestInMemoryRepo_getAll(t *testing.T) {
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
			name: "customers present in database",
			fields: fields{
				customers: []Customer{
					{
						Id: "hs",
						CustomerDetails: CustomerDetails{
							Name:      "hardik",
							Address:   "udaipur",
							ContactNo: 9649127550,
						},
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127550,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "empty database",
			fields: fields{
				customers: []Customer{},
			},
			wantCustomers: []Customer{},
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}

			gotCustomers, gotErr := repo.getAll()

			if gotErr != nil {
				t.Errorf("got an error :%q", gotErr)
			}

			if !reflect.DeepEqual(gotCustomers, tt.wantCustomers) {
				t.Errorf("customers list should be\nwant customers %+v\nbut got %+v", tt.wantCustomers, gotCustomers)
			}
		})
	}
}

func TestInMemoryRepo_getById(t *testing.T) {
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
			name: "getting existing customer",
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
					{
						Id: "mm",
						CustomerDetails: CustomerDetails{
							Name:      "vs",
							Address:   "gj",
							ContactNo: 9649127559,
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
					ContactNo: 9649127559,
				},
			},
			wantErr: nil,
		},
		{
			name: "getting customer that does not exist",
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
					{
						Id: "mm",
						CustomerDetails: CustomerDetails{
							Name:      "vs",
							Address:   "gj",
							ContactNo: 9649127559,
						},
					},
				},
			},
			args: args{
				id: "hm",
			},
			wantCustomer: Customer{},
			wantErr:      ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &InMemoryRepo{customers: tt.fields.customers}

			gotCustomer, gotErr := repo.getById(tt.args.id)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("got error %q but want error :%q", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotCustomer, tt.wantCustomer) {
				t.Errorf("got %+v\n want %+v", gotCustomer, tt.wantCustomer)
			}
		})
	}
}
