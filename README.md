# gFly Database

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

Fluent Model - flexible and powerful Data-Access Layer. Build on top of [Fluent SQL](https://github.com/JiveIO/FluentSQL)

### Usage

Install
```bash
go get -u github.com/gflydev/db@latest

# PostgreSQL
go get -u github.com/gflydev/db/psql@latest

# MySQL
go get -u github.com/gflydev/db/mysql@latest
```

Quick usage `main.go`
```go
import (
    mb "github.com/gflydev/db" 
    // PostgreSQL
    _ "github.com/gflydev/db/psql"
    // MySQL
    //_ "github.com/gflydev/db/mysql"
)

func main() {
    // Load Model builder
    mb.Load()
}

func queryUser() {
    try.Perform(func() {
        dbInstance := mb.Instance()
        if dbInstance == nil {
            panic("Database Model is NULL")
        }

        // Defer a rollback in case anything fails.
        defer func(db *mb.DBModel) {
            _ = db.Rollback()
        }(dbInstance)

        var user models.User
        err := dbInstance.First(&user)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("User %v\n", user)
    }).Catch(func(e try.E) {
        log.Printf("Error %v\n", e)
    })
}
```

### Generic DAO

Basic methods to create CRUD actions 

```go
import (
    mb "github.com/gflydev/db" // Model builder
    qb "github.com/jivegroup/fluentsql" // Query builder
)
...
// ----- GetModelByID -----
user, err := mb.GetModelByID[models.User](1)
if err != nil {
    log.Fatal(err)
}
log.Info("Get \n", user1.Email)

// ----- CreateModel -----
err = mb.CreateModel(&models.User{
    Email:        "john@gmail.com",
    Password:     "02j33ih32i3",
    Fullname:     "John Man",
    Phone:        "0989712353",
    Token:        sql.NullString{},
    Status:       "active",
    CreatedAt:    time.Time{},
    Avatar:       sql.NullString{},
    UpdatedAt:    time.Time{},
    VerifiedAt:   sql.NullTime{},
    BlockedAt:    sql.NullTime{},
    DeletedAt:    sql.NullTime{},
    LastAccessAt: sql.NullTime{},
})
if err != nil {
    log.Fatal(err)
}

// ----- FindModels -----
users, total, err := mb.FindModels[models.User](1, 100, "id", mb.Desc, mb.Condition{
        Field: "id",
        Opt:   mb.NotEq,
        Value: 0,
    })
if err != nil {
    log.Fatal(err)
}
log.Info("Find \n", total)
for _, user := range users {
    log.Info("User\n", user.Email)
}

// ----- UpdateModel -----
user1.Fullname = "Admin"
if err := mb.UpdateModel(user1); err != nil {
    log.Fatal(err)
}
log.Info("Update \n", user1.Fullname)
```

### More using `FluentSQL` and `FluentModel`
```go
import (
    mb "github.com/gflydev/db" // Model builder
    qb "github.com/jivegroup/fluentsql" // Query builder
)
```

### User `model`
We have model `User` is `struct` type. Every `fluent model` must have `MetaData` field to specify some table metadata 
to link from `model` to `table`.
```go
// User model
type User struct {
    // Table meta data
    MetaData fm.MetaData `db:"-" model:"table:users"`

    // Table fields
    Id   int            `db:"id" model:"type:serial,primary"`
    Name sql.NullString `db:"name" model:"type:varchar(255)"`
    Age  uint8          `db:"age" model:"type:numeric"`

	// Extra fields
    TotalAge int `db:"total_age" model:"type:numeric"`
}
```

### Fluent model object

A new ModelBuilder instance to help you perform `query`, `create`, `update` or `delete` data.
```go
db := mb.Instance()
if db == nil {
    panic("Database Model is NULL")
}

// Defer a rollback in case anything fails.
defer func(db *mb.DBModel) {
    _ = db.Rollback()
}(db)

var err error
```

## Query data

**Get the first ORDER BY id ASC**
```go
var user User
err = db.First(&user)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user)
```

**Get the last ORDER BY id DESC**
```go
var user1 User
err = db.Last(&user1)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user1)
```

**Get a random item**
```go
// Get a random item
var user2 User
err = db.Take(&user2)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user2)
```

**Get first by ID**
```go
var user3 User
err = db.First(&user3, 103)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user3)
```

**Get first by model ID**
```go
var user4 User
user4 = User{Id: 103}
err = db.First(&user4)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user4)
```

**Get first by model**
```go
var user5 User
err = db.Model(User{Id: 102}).
	First(&user5)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user5)
```

**Get first by Where condition**
```go
var user6 User
err = db.Where("Id", Eq, 100).
	First(&user6)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user6)
```

**Get first by WhereGroup**
```go
var user7 User
err = db.Where("Id", Eq, 100).
    WhereGroup(func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("age", Eq, 42).
            WhereOr("age", Eq, 39)

        return &query
    }).First(&user7)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user7)
```

**Query all**
```go
var users []User
_, err = db.Find(&users)
if err != nil {
    log.Fatal(err)
}
for _, user := range users {
    log.Printf("User %v\n", user)
}
```

**Query by model condition**
```go
var users1 []User
_, err = db.Model(&User{Age: 20}).Find(&users1)
if err != nil {
    log.Fatal(err)
}
for _, user := range users1 {
    log.Printf("User %v\n", user)
}
```

**Query by Slice (list of ID)**
```go
var users2 []User
_, err = db.Find(&users2, []int{144, 145, 146})
if err != nil {
    log.Fatal(err)
}
for _, user := range users2 {
    log.Printf("User %v\n", user)
}
```

**Query with SELECT|GROUP BY|HAVING**
```go
var users3 []UserTotal
_, err = db.Model(&UserTotal{}).
    Select("name, sum(age) as total_age").
    GroupBy("name").
    Having("name", Eq, "vinh").
    Find(&users3)
if err != nil {
    log.Fatal(err)
}
for _, user := range users3 {
    log.Printf("User %v\n", user)
}
```

**Query with JOIN**
```go
var users4 []UserJoin
_, err = db.Model(&UserJoin{}).
    Select("name, age, email, phone").
    Join(mb.InnerJoin, "user_details", mb.Condition{
        Field: "users.id",
        Opt:   Eq,
        Value: mb.ValueField("user_details.user_id"),
    }).
    Where("users.name", Eq, "Kite").
    Find(&users4)
if err != nil {
    log.Fatal(err)
}
for _, user := range users4 {
    log.Printf("User %v\n", user)
}
```

**Query with raw SQL**
```go
var users5 []User
_, err = db.Raw("SELECT * FROM users WHERE name = ?", "Kite").
    Find(&users5)
if err != nil {
    log.Fatal(err)
}
for _, user := range users5 {
    log.Printf("User %v\n", user)
}
```

**Query with paging info**
```go
var (
    users6 []User
    total  int
)
total, err = db.Limit(10, 0).Find(&users6)

if err != nil {
    log.Fatal(err)
}
log.Printf("Total %d\n", total)
for _, user := range users6 {
    log.Printf("User %v\n", user)
}
```

## Create data

**Create from a model**
```go
user := User{
    Name: sql.NullString{String: "Vinh", Valid: true},
    Age:  42,
}
// Create new row into table `users`
err = db.Create(&user) 
if err != nil {
    log.Fatal(err)
}
log.Printf("User ID: %d", user.Id)
```

**Create from model - Omit a column**
```go
userDetail := UserDetail{
    UserId: 1,
    Email:  "vinh@mail.com",
    Phone:  1989831911,
}

// Create new row but skips data of column `phone`
err = db.Omit("phone").
	Create(&userDetail)
if err != nil {
    log.Fatal(err)
}
log.Printf("User detail ID: %d", userDetail.Id)
```

**Create from model - For some specific columns**
```go
userDetail = UserDetail{
    UserId: 1,
    Email:  "vinh.vo@gmail.com",
    Phone:  975821086,
}

// Create new row but only data for column `user_id` and `email`
err = db.Select("user_id", "email").
	Create(&userDetail)
if err != nil {
    log.Fatal(err)
}
log.Printf("User detail ID: %d", userDetail.Id)
```

**Create from Slice models**
```go
var users []*User
users = append(users, &User{
    Name: sql.NullString{String: "John", Valid: true},
    Age:  39,
})
users = append(users, &User{
    Name: sql.NullString{String: "Kite", Valid: true},
    Age:  42,
})
err = db.Create(users)
if err != nil {
    log.Fatal(err)
}

for _, user := range users {
    log.Printf("User ID: %d", user.Id)
}
```

**Create from Map column keys**
```go
user = User{}
err = db.Model(&user).
	Create(map[string]interface{}{
        "Name": "John Lyy",
        "Age":  39,
    })
if err != nil {
    log.Fatal(err)
}
log.Printf("User ID: %d", user.Id)
```

## Update data

**Update by model**
```go
var user User
err = db.First(&user)
user.Name = sql.NullString{
    String: "Cat John",
    Valid:  true,
}

err = db.Update(user)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user)
```

**Update by model and condition**
```go
var user1 User
err = db.First(&user1)
user1.Name = sql.NullString{
    String: "Cat John",
    Valid:  true,
}
user1.Age = 100

err = db.
    Where("id", Eq, 1).
    Update(user1)

if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user1)
```

**Update by Map**
```go
var user2 User
err = db.First(&user2)
err = db.Model(&user2).
    Omit("Name").
    Update(map[string]interface{}{"Name": "Tah Go Tab x3", "Age": 88})

if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user2)
```

## Delete data

**Delete by Model**
```go
var user User
err = db.First(&user)
err = db.Delete(user)
if err != nil {
    log.Fatal(err)
}
```

**Delete by ID**
```go
err = db.Delete(User{}, 157)
if err != nil {
    log.Fatal(err)
}
```

**Delete by List ID**
```go
err = db.Delete(User{}, []int{154, 155, 156})
if err != nil {
    log.Fatal(err)
}
```

**Delete by Where condition**
```go
err = db.Where("Id", Eq, 153).
	Delete(&User{})
if err != nil {
    log.Fatal(err)
}
```

## RAW SQLs

```go
// -------- Insert --------
var user User
err = db.Raw("INSERT INTO users(name, age) VALUES($1, $2)", "Kite", 43).
    Create(&user)
if err != nil {
    log.Fatal(err)
}
log.Printf("User %v\n", user)

// -------- Update --------
err = db.Raw("UPDATE users SET name = $1, age = $2 WHERE id= $3", "Kite - Tola", 34, 1).
    Update(&User{})
if err != nil {
    log.Fatal(err)
}

// -------- Get One --------
var user2 User
err = db.Raw("SELECT * FROM users WHERE id=$1", 1).
    First(&user2)
log.Printf("User %v\n", user2)

// -------- Select --------
var userList []User
var total int
total, err = db.Raw("SELECT * FROM users").
    Find(&userList)
log.Printf("Total %v\n", total)

for _, _user := range userList {
    log.Printf("User %v\n", _user)
}

// -------- Delete --------
err = db.Raw("DELETE FROM users WHERE id > $1", 1).
    Delete(&User{})
if err != nil {
    log.Fatal(err)
}
```
