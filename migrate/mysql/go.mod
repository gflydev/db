module github.com/gflydev/db/migrate/mysql

go 1.26.0

replace github.com/gflydev/db/migrate => ../

require (
	github.com/gflydev/core v1.18.3
	github.com/gflydev/db/migrate v0.0.0
	github.com/go-sql-driver/mysql v1.10.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/molecule-man/go-brrr v1.1.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.74.0 // indirect
	golang.org/x/crypto v0.57.0 // indirect
)
