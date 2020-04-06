package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	//static files
	e.Static("/", "web")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/dg/update", updateStatus)
	e.PUT("/dg/update", updateStatus)
	e.GET("/dg/status/:id", getStatus)
	e.GET("/dg/status/list", getStatusList)
	e.GET("/dg/list/:serviceitem", getUniq)
	e.GET("/dg/servicefilter/:item/:itemval", serviceFilter)
	e.DELETE("/dg/service/:id", deleteService)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
