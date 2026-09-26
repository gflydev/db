module github.com/gflydev/db/migrate/mysql

go 1.24.0

replace github.com/gflydev/db/migrate => ../

require (
	github.com/gflydev/core v1.17.13
	github.com/gflydev/db/migrate v0.0.0
	github.com/go-sql-driver/mysql v1.9.3
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.67.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
)
