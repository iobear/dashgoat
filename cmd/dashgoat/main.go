package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config Configer
var dashgoat_ready bool
var updatekey string
var ss Services
var bb Backlog
var dashName string

func main() {
	var ipport string
	var webpath string
	var weblog string
	var configfile string
	var oneBuddyUrl string
	var oneBuddyKey string
	var oneBuddyName string

	ss.serviceStateList = make(map[string]ServiceState)
	bb.buddyBacklog = make(map[string][]string)
	bb.StateDown = make(map[string]int64)

	e := echo.New()

	flag.StringVar(&ipport, "ipport", ":1323", "Specify <ip>:<port>")
	flag.StringVar(&weblog, "weblog", "off", "HTTP log <on/off>")
	flag.StringVar(&webpath, "webpath", "/", "Specify added url http://host:port/<path> Default: /")
	flag.StringVar(&updatekey, "updatekey", "changeme", "Specify key to API update")
	flag.StringVar(&dashName, "dashname", "", "Dashboard name")
	flag.StringVar(&configfile, "configfile", "dashgoat.yaml", "Name of configfile")
	flag.StringVar(&oneBuddyUrl, "buddyurl", "", "Buddy url")
	flag.StringVar(&oneBuddyKey, "buddykey", "", "Buddy update key, empty for same key")
	flag.StringVar(&oneBuddyName, "buddyname", "", "Buddy name")
	flag.Parse()

	err := config.InitConfig(configfile)
	if err != nil {
		fmt.Println(err)
	}

	if oneBuddyUrl != "" {
		if oneBuddyKey == "" {
			oneBuddyKey = updatekey
		}
		config.OneBuddy(oneBuddyUrl, oneBuddyKey, oneBuddyName)
	}

	pathStartsWith := strings.HasPrefix(webpath, "/")
	if !pathStartsWith {
		webpath = "/" + webpath
	}

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
	e.GET(add2url(webpath, "/status/listmso"), getStatusListMSO)
	e.GET(add2url(webpath, "/list/:serviceitem"), getUniq)
	e.GET(add2url(webpath, "/servicefilter/:item/:itemval"), serviceFilter)
	e.DELETE(add2url(webpath, "/service/:id"), deleteService)
	e.GET(add2url(webpath, "/health"), health)

	fmt.Println("Starting dashGoat.. Dashboard name: " + dashName + " ")

	go lostProbeTimer()
	go findBuddy()

	// Start server
	e.Logger.Fatal(e.Start(ipport))

}
