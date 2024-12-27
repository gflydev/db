module github.com/gflydev/db/psql

go 1.22.6

require github.com/jmoiron/sqlx v1.4.0

replace github.com/gflydev/db => ../

require (
	github.com/gflydev/db v1.2.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/jiveio/fluentsql v1.4.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/gflydev/core v1.10.6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
