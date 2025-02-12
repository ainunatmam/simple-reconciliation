package routes

import (
	"time"

	"fiber.misetaku.net/services"
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	F *fiber.App
}

func (r Route) InitRoute() {
	r.F.Get("/", func(c *fiber.Ctx) error {
		start,_ := time.Parse("2006-01-02", c.Query("start_date"))
		end,_ := time.Parse("2006-01-02", c.Query("end_date"))
		services := services.ManualRecon{
			TransactionPath: "./files/system_transaction.csv",
			StatementPath: "./files/statements",
			StartDate: start,
			EndDate: end,
		}
		result := services.Call()
		return c.JSON(result)
	})
}
