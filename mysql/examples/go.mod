module examples

go 1.24.0

toolchain go1.24.2

replace (
	github.com/gflydev/db => ../../
	github.com/gflydev/db/mysql => ../
)

require github.com/gflydev/core v1.17.7

require github.com/gflydev/view/pongo v1.0.3

require (
	github.com/gflydev/db v1.11.0
	github.com/gflydev/db/v2 v2.0.1
)

require (
	github.com/gflydev/db/mysql v1.4.9
	github.com/gflydev/session v1.0.3
	github.com/gflydev/session/memory v1.0.3
	github.com/joho/godotenv v1.5.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/flosch/pongo2/v6 v6.0.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/jivegroup/fluentsql v1.5.4 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.66.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
