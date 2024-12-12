package main

import (
	"examples/models"
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	dbPSQL "github.com/gflydev/db/psql"
	"github.com/gflydev/session"
	sessionMemory "github.com/gflydev/session/memory"
	"github.com/gflydev/view/pongo"

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
