module github.com/gflydev/db/psql

go 1.22.6

replace github.com/gflydev/db => ../

require github.com/jmoiron/sqlx v1.4.0

require github.com/gflydev/db v0.0.0-00010101000000-000000000000

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)
