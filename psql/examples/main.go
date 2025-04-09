package main

import (
	"database/sql"
	"examples/models"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	dbPSQL "github.com/gflydev/db/psql"
	"github.com/gflydev/session"
	sessionMemory "github.com/gflydev/session/memory"
	"github.com/gflydev/view/pongo"
	qb "github.com/jiveio/fluentsql"
	"time"

	// Autoload .env file
	_ "github.com/joho/godotenv/autoload"
)

// =========================================================================================
//                                     Default API
// =========================================================================================

// NewDefaultApi As a constructor to create new API.
func NewDefaultApi() *DefaultApi {
	return &DefaultApi{}
}

// DefaultApi API struct.
type DefaultApi struct {
	core.Api
}

func (h *DefaultApi) Handle(c *core.Ctx) error {
	return c.JSON(core.Data{
		"name":   core.AppName,
		"server": core.AppURL,
	})
}

// =========================================================================================
//                                     Home page
// =========================================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	core.Page
}

func (m *HomePage) Handle(c *core.Ctx) error {
	c.SetSession("foo", utils.UnsafeStr(utils.RandByte(make([]byte, 128))))

	// Database
	queryUser()

	// Generic DAO
	genericDao()

	return c.View("home", core.Data{
		"title": "gFly | Laravel inspired web framework written in Go",
	})
}

// =========================================================================================
//                                     Session page
// =========================================================================================

// NewSessionPage As a constructor to create a Session Page.
func NewSessionPage() *SessionPage {
	return &SessionPage{}
}

type SessionPage struct {
	core.Page
}

func (m *SessionPage) Handle(c *core.Ctx) error {
	foo := c.GetSession("foo")

	return c.View("session", core.Data{
		"title": "gFly | Laravel inspired web framework written in Go",
		"foo":   foo,
	})
}

// =========================================================================================
//                                     Routers
// =========================================================================================

func router(g core.IFly) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	// API Routers
	g.Group(prefixAPI, func(apiRouter *core.Group) {
		apiRouter.GET("/info", NewDefaultApi())
	})

	// Web Routers
	g.GET("/home", NewHomePage())
	g.GET("/session", NewSessionPage())
}

// =========================================================================================
//                                     Application
// =========================================================================================

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
		log.Info("User\n", user)
	}).Catch(func(e try.E) {
		log.Error("Error\n", e)
	})
}

func genericDao() {
	try.Perform(func() {
		// ----- GetModelByID -----
		user1, err := mb.GetModelByID[models.User](1)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Get \n", user1.Email)

		// ----- CreateModel
		err = mb.CreateModel(&models.User{
			Email:        gofakeit.Email(),
			Password:     gofakeit.Password(true, true, true, true, true, 6),
			Fullname:     gofakeit.Name(),
			Phone:        gofakeit.Phone(),
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
		users, total, err := mb.FindModels[models.User](1, 100, "id", qb.Desc, qb.Condition{
			Field: "id",
			Opt:   qb.NotEq,
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
	}).Catch(func(e try.E) {
		log.Error("Error\n", e)
	})
}

func main() {
	app := core.New()

	// Register view
	core.RegisterView(pongo.New())

	// Setup session
	session.Register(sessionMemory.New())
	core.RegisterSession(session.New())

	// Register DB driver & Load Model builder
	mb.Register(dbPSQL.New())
	mb.Load()

	// Register router
	app.RegisterRouter(router)

	app.Run()
}
