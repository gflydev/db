module github.com/gflydev/db/mysql

go 1.22.6

replace github.com/gflydev/db => ../

require github.com/jmoiron/sqlx v1.4.0
require github.com/gflydev/db v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
)
