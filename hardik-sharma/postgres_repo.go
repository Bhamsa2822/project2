package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type postgresRepo struct {
	db *bun.DB
}

func NewPostgresRepo(db *bun.DB) *postgresRepo {
	return &postgresRepo{
		db: db,
	}
}

func (repo *postgresRepo) create(customer Customer) error {
	if _, err := repo.db.NewInsert().Model(&customer).Exec(context.Background()); err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) && pgdriverErr.IntegrityViolation() {
			return ErrConflict
		}

		return err
	}

	return nil
}

func (repo *postgresRepo) getAll() ([]Customer, error) {
	customers := []Customer{}
	if err := repo.db.NewSelect().Model(&customers).Scan(context.Background()); err != nil {
		return customers, err
	}

	return customers, nil
}

func (repo *postgresRepo) getById(id string) (Customer, error) {
	var customer Customer
	if err := repo.db.NewSelect().Model(&customer).Where("id = ?", id).Scan(context.Background()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Customer{}, ErrNotFound
		}

		return Customer{}, err
	}

	return customer, nil
}

func (repo *postgresRepo) update(id string, customer Customer) error {
	res, err := repo.db.NewUpdate().Model(&customer).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return err
	}

	rowsAffectCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffectCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (repo *postgresRepo) delete(id string) error {
	res, err := repo.db.NewDelete().Model((*Customer)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return err
	}

	rowAffectCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffectCount == 0 {
		return ErrNotFound
	}

	return nil
}
