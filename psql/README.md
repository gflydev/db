# gFly Database - PostgreSQL

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

Fluent Model - flexible and powerful Data-Access Layer. Build on top of [Fluent SQL](https://github.com/JiveIO/FluentSQL)

### Usage

Install
```bash
go get -u github.com/gflydev/db@v1.0.0
go get -u github.com/gflydev/db/psql@v1.1.0
```

Quick usage `main.go`
```go
import (
    mb "github.com/gflydev/db"
    dbPSQL "github.com/gflydev/db/psql"
)

func main() {
    // Register DB driver & Load Model builder
    mb.Register(dbPSQL.New())
    mb.Load()
}
```