package main

import (
	"flag"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var updatekey string
var ss Services

func main() {
	var ipport string
	var webpath string
	var weblog string

	ss.serviceStateList = make(map[string]ServiceState, 0)

	flag.StringVar(&ipport, "ipport", ":1323", "Specify <ip>:<port>")
	flag.StringVar(&weblog, "weblog", "off", "HTTP log <on/off>")
	flag.StringVar(&webpath, "webpath", "/", "Specify added url http://host:port/<path> Default: /")
	flag.StringVar(&updatekey, "updatekey", "changeme", "Specify key to API update")
	flag.Parse()

	pathStartsWith := strings.HasPrefix(webpath, "/")
	if pathStartsWith == false {
		webpath = "/" + webpath
	}

	e := echo.New()

	e.HideBanner = true

	//static files
	e.Static(webpath, "web")

	if weblog == "on" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	e.POST(add2url(webpath, "/update"), updateStatus)
	e.GET(add2url(webpath, "/status/:id"), getStatus)
	e.GET(add2url(webpath, "/status/list"), getStatusList)
	e.GET(add2url(webpath, "/list/:serviceitem"), getUniq)
	e.GET(add2url(webpath, "/servicefilter/:item/:itemval"), serviceFilter)
	e.DELETE(add2url(webpath, "/service/:id"), deleteService)
	e.GET(add2url(webpath, "/health"), health)

	print("dashGoat ")

	// Start server
	e.Logger.Fatal(e.Start(ipport))
}
