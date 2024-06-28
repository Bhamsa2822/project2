# Start postgres db

docker run --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name customers_db postgres


url: localhost:5432
username: postgres
password: postgres
database: postgres