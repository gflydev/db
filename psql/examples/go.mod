module examples

go 1.22.6

replace github.com/gflydev/db => ../../

replace github.com/gflydev/db/psql => ../

require github.com/gflydev/core v1.10.6

require github.com/gflydev/view/pongo v1.0.2

require github.com/gflydev/db/psql v1.1.0

require github.com/gflydev/db v1.2.0

require (
	github.com/gflydev/session v1.0.1
	github.com/gflydev/session/memory v1.0.1
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/flosch/pongo2/v6 v6.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jiveio/fluentsql v1.4.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
