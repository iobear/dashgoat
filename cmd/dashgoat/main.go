package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config Configer
var buddyCli Buddy
var dashgoat_ready bool
var ss Services
var bb Backlog

func main() {
	var configfile string

	ss.serviceStateList = make(map[string]ServiceState)
	bb.buddyBacklog = make(map[string][]string)
	bb.StateDown = make(map[string]int64)

	e := echo.New()

	flag.StringVar(&config.IPport, "ipport", ":2000", "Specify <ip>:<port>")
	flag.StringVar(&config.WebLog, "weblog", "off", "HTTP log <on/off>")
	flag.StringVar(&config.WebPath, "webpath", "/", "Specify added url http://host:port/<path> Default: /")
	flag.StringVar(&config.UpdateKey, "updatekey", "changeme", "Specify key to API update")
	flag.StringVar(&config.DashName, "dashname", "", "Dashboard name")
	flag.StringVar(&configfile, "configfile", "dashgoat.yaml", "Name of configfile")
	flag.StringVar(&buddyCli.Url, "buddyurl", "", "Buddy url")
	flag.StringVar(&buddyCli.Key, "buddykey", "", "Buddy update key, empty for same key")
	flag.StringVar(&buddyCli.Name, "buddyname", "", "Buddy name")
	flag.Parse()

	err := config.InitConfig(configfile)
	if err != nil {
		fmt.Println(err)
	}

	pathStartsWith := strings.HasPrefix(config.WebPath, "/")
	if !pathStartsWith {
		config.WebPath = "/" + config.WebPath
	}

	e.HideBanner = true

	//static files
	e.Static(config.WebPath, "web")

	if config.WebLog == "on" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	e.POST(add2url(config.WebPath, "/update"), updateStatus)
	e.GET(add2url(config.WebPath, "/status/:id"), getStatus)
	e.GET(add2url(config.WebPath, "/status/list"), getStatusList)
	e.GET(add2url(config.WebPath, "/status/listmso"), getStatusListMSO)
	e.GET(add2url(config.WebPath, "/list/:serviceitem"), getUniq)
	e.GET(add2url(config.WebPath, "/servicefilter/:item/:itemval"), serviceFilter)
	e.DELETE(add2url(config.WebPath, "/service/:id"), deleteService)
	e.GET(add2url(config.WebPath, "/health"), health)

	fmt.Println("Starting dashGoat.. Dashboard name: " + config.DashName + " ")

	go lostProbeTimer()
	go findBuddy()

	// Start server
	e.Logger.Fatal(e.Start(config.IPport))

}
