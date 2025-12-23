module examples

go 1.24.0

replace (
	github.com/gflydev/db => ../../
	github.com/gflydev/db/psql => ../
)

require github.com/gflydev/core v1.17.13

require github.com/gflydev/view/pongo v1.0.3

require github.com/gflydev/db/psql v1.4.9

require github.com/gflydev/db v1.11.0

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0
	github.com/gflydev/session v1.0.3
	github.com/gflydev/session/memory v1.0.3
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/flosch/pongo2/v6 v6.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jivegroup/fluentsql v1.5.4 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.68.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
