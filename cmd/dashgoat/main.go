package main

import (
	"flag"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var ipport string
	var webpath string
	var weblog string

	flag.StringVar(&ipport, "ipport", ":1323", "Specify <ip>:<port>")
	flag.StringVar(&weblog, "weblog", "off", "Specify <ip>:<port>")
	flag.StringVar(&webpath, "webpath", "/", "Specify added url http://host:port/<path> Default: /")
	flag.Parse()

	pathStartsWith := strings.HasPrefix(webpath, "/")
	if pathStartsWith == false {
		webpath = "/" + webpath
	}

	e := echo.New()

	//static files
	e.Static(webpath, "web")

	// Middleware
	if weblog == "on" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	// Routes
	e.POST(add2url(webpath, "/update"), updateStatus)
	e.PUT(add2url(webpath, "/update"), updateStatus)
	e.GET(add2url(webpath, "/status/:id"), getStatus)
	e.GET(add2url(webpath, "/status/list"), getStatusList)
	e.GET(add2url(webpath, "/list/:serviceitem"), getUniq)
	e.GET(add2url(webpath, "/servicefilter/:item/:itemval"), serviceFilter)
	e.DELETE(add2url(webpath, "/service/:id"), deleteService)

	// Start server
	e.Logger.Fatal(e.Start(ipport))
}
