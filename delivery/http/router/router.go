package router

import (
	"database/sql"
	"digibank/delivery/http/controllers"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, db *sql.DB, c *fiber.Ctx) {
	app.Post("/api/user/register", func(c *fiber.Ctx) error {
		return controllers.Register(db, c)
	})
	app.Post("/api/user/login", func(c *fiber.Ctx) error {
		return controllers.Login(db, c)
	})
	app.Put("/api/user/update", func(c *fiber.Ctx) error {
		return controllers.UpdateUser(db, c)
	})
	app.Delete("/api/user/delete", func(c *fiber.Ctx) error {
		return controllers.DeactiveUserAccount(db, c)
	})
	app.Get("/api/user/info", func(c *fiber.Ctx) error {
		return controllers.GetUser(db, c)
	})
	app.Post("/api/account/check-pin", func(c *fiber.Ctx) error {
		return controllers.CheckPin(db, c)
	})
	app.Post("/api/account/create-account", func(c *fiber.Ctx) error {
		return controllers.CreateAccount(db, c)
	})
	app.Get("/api/account/all", func(c *fiber.Ctx) error {
		return controllers.GetAllAcct(db, c)
	})

	app.Get("/api/account/:accountno", func(c *fiber.Ctx) error {
		return controllers.GetAcct(db, c)
	})
	app.Delete("/api/account/:accountno", func(c *fiber.Ctx) error {
		return controllers.CloseAcct(db, c)
	})
	app.Post("/api/transaction/topup", func(c *fiber.Ctx) error {
		return controllers.AddBalance(db, c)
	})
	app.Post("/api/transaction/overbook", func(c *fiber.Ctx) error {
		return controllers.Overbooking(db, c)
	})
	app.Get("/api/transaction/mutation/:accountnumber", func(c *fiber.Ctx) error {
		return controllers.GetUserMutation(db, c)
	})
	app.Post("/api/transaction/time-depo-sim", func(c *fiber.Ctx) error {
		return controllers.TimeDepositSimulation(db, c)
	})
	app.Get("/api/welcome", func(c *fiber.Ctx) error {
		return controllers.Welcome(db, c)
	})
}
