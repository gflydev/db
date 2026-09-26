module github.com/gflydev/db/migrate/postgres

go 1.24.0

replace github.com/gflydev/db/migrate => ../

require (
	github.com/gflydev/core v1.17.13
	github.com/gflydev/db/migrate v0.0.0
	github.com/jackc/pgx/v5 v5.7.5
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.67.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)
