package main

import (
	"billing_engine/config"
	"billing_engine/repository"
	"billing_engine/scheduler"
	"billing_engine/services/billing"
	"billing_engine/utils"
	"billing_engine/utils/database"
	"billing_engine/utils/route"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	utils.AddValidator(e)
	conf := config.NewConfig()
	db := database.NewDatabase(conf)
	repository := repository.NewRepository(db)
	billingService := billing.NewService(repository)
	route.NewRoute(e, route.Services{
		BillingService: billingService,
	})
	scheduler.NewScheduler(billingService)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
}
