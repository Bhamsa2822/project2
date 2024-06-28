-- +goose Up

CREATE TABLE customers(
   id TEXT PRIMARY KEY,
   customerdetails_name TEXT,
   customerdetails_address TEXT,
   customerdetails_contact_no BIGINT
);

-- +goose Down
DROP TABLE customers;