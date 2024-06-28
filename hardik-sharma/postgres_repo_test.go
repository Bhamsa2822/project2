package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func setupDB(t *testing.T, existingCustomers []Customer) *bun.DB {
	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		t.Fatal("db connection check failed:", err)
	}

	if _, err := db.Query("TRUNCATE TABLE customers"); err != nil {
		t.Fatal("failed to truncate table:", err)
	}

	if len(existingCustomers) > 0 {
		if _, err := db.NewInsert().Model(&existingCustomers).Exec(context.Background()); err != nil {
			t.Fatal("failed to add customers:", err)
		}
	}

	return db
}

func Test_postgresRepo_create(t *testing.T) {
	type args struct {
		customer Customer
	}

	tests := []struct {
		name              string
		existingCustomers []Customer
		args              args
		wantCustomers     []Customer
		wantErr           error
	}{
		{
			name: "adding new customer",
			existingCustomers: []Customer{
				{
					Id: "ht",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			args: args{
				customer: Customer{
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
					Id: "ht",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
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
			wantErr: nil,
		},
		{
			name: "adding existing customer",
			existingCustomers: []Customer{
				{
					Id: "ht",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			args: args{
				customer: Customer{
					Id: "ht",
					CustomerDetails: CustomerDetails{
						Name:      "h",
						Address:   "u",
						ContactNo: 7777777777,
					},
				},
			},
			wantCustomers: []Customer{
				{
					Id: "ht",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9999999999,
					},
				},
			},
			wantErr: ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t, tt.existingCustomers)
			repo := NewPostgresRepo(db)

			gotErr := repo.create(tt.args.customer)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to be same")

			var gotCustomers []Customer
			if err := db.NewSelect().Model(&gotCustomers).Scan(context.Background()); err != nil {
				t.Fatal("failed to fetch customers", err)
			}

			assert.Equal(t, tt.wantCustomers, gotCustomers, "expected customers to be same")
		})
	}
}

func Test_postgresRepo_update(t *testing.T) {
	type args struct {
		id             string
		updateCustomer Customer
	}

	tests := []struct {
		name              string
		existingCustomers []Customer
		args              args
		wantCustomers     []Customer
		wantErr           error
	}{
		{
			name: "updating existing customer",
			existingCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "uadipur",
						ContactNo: 9649127559,
					},
				},
				{
					Id: "ps",
					CustomerDetails: CustomerDetails{
						Name:      "parmavrr",
						Address:   "uadipur",
						ContactNo: 9649127559,
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
					Id: "ps",
					CustomerDetails: CustomerDetails{
						Name:      "parmavrr",
						Address:   "uadipur",
						ContactNo: 9649127559,
					},
				},
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "Ahmedabad",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "updating non existing customer",
			existingCustomers: []Customer{
				{
					Id: "hk",
					CustomerDetails: CustomerDetails{
						Name:      "h",
						Address:   "raj",
						ContactNo: 9649127559,
					},
				},
				{
					Id: "ps",
					CustomerDetails: CustomerDetails{
						Name:      "parmavrr",
						Address:   "uadipur",
						ContactNo: 9649127559,
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
				{
					Id: "ps",
					CustomerDetails: CustomerDetails{
						Name:      "parmavrr",
						Address:   "uadipur",
						ContactNo: 9649127559,
					},
				},
			},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t, tt.existingCustomers)
			repo := NewPostgresRepo(db)

			gotErr := repo.update(tt.args.id, tt.args.updateCustomer)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to be same")

			var customers []Customer
			if err := db.NewSelect().Model(&customers).Scan(context.Background()); err != nil {
				t.Fatal("unable to fetch customers list :", err)
			}

			assert.Equal(t, tt.wantCustomers, customers, "expect customers to be same")
		})
	}
}

func Test_postgresRepo_getAll(t *testing.T) {
	tests := []struct {
		name              string
		existingCustomers []Customer
		wantCustomers     []Customer
		wantErr           error
	}{
		{
			name: "customers present in database",
			existingCustomers: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127550,
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
			name:              "empty database",
			existingCustomers: []Customer{},
			wantCustomers:     []Customer{},
			wantErr:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t, tt.existingCustomers)
			repo := NewPostgresRepo(db)

			gotCustomers, gotErr := repo.getAll()

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to be same")

			assert.Equal(t, tt.wantCustomers, gotCustomers, "expect customer list to be matched")
		})
	}
}

func Test_postgresRepo_getById(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name              string
		existingCustomers []Customer
		args              args
		wantCustomer      Customer
		wantErr           error
	}{
		{
			name: "getting existing customer",
			existingCustomers: []Customer{
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
			existingCustomers: []Customer{
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
			args: args{
				id: "hm",
			},
			wantCustomer: Customer{},
			wantErr:      ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t, tt.existingCustomers)
			repo := NewPostgresRepo(db)

			gotCustomer, gotErr := repo.getById(tt.args.id)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to be same")

			assert.Equal(t, tt.wantCustomer, gotCustomer, "expect customer to be matched")
		})
	}
}

func Test_postgresRepo_delete(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name             string
		existingCustomer []Customer
		args             args
		wantCustomers    []Customer
		wantErr          error
	}{
		{
			name: "deleting existing customer",
			existingCustomer: []Customer{
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
			existingCustomer: []Customer{
				{
					Id: "hs",
					CustomerDetails: CustomerDetails{
						Name:      "hardik",
						Address:   "udaipur",
						ContactNo: 9649127559,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t, tt.existingCustomer)
			repo := NewPostgresRepo(db)

			gotErr := repo.delete(tt.args.id)

			assert.ErrorIs(t, tt.wantErr, gotErr, "expect error to be same")

			var customerList []Customer
			if err := db.NewSelect().Model(&customerList).Scan(context.Background()); err != nil {
				t.Fatal("failed to fetch data :", err)
			}

			assert.Equal(t, tt.wantCustomers, customerList, "expect customer list to be same")
		})
	}
}
