/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"

	dg "github.com/iobear/dashgoat/common"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	//ServiceState for MSO output
	ServiceStateMSO struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	Services struct {
		mutex            sync.RWMutex
		serviceStateList map[string]dg.ServiceState
	}
)

//go:embed web
var embededFiles embed.FS

var config Configer
var buddy_cli Buddy
var nsconfig string
var ss Services

var backlog Backlog
var serviceStateCollector *ServiceStateCollector

func main() {
	var configfile string

	ss.serviceStateList = make(map[string]dg.ServiceState)
	backlog.buddyBacklog = make(map[string][]string)
	backlog.StateDown = make(map[string]int64)

	e := echo.New()

	flag.StringVar(&config.IPport, "ipport", ":2000", "Specify <ip>:<port>")
	flag.StringVar(&config.WebLog, "weblog", "off", "HTTP log <on/off>")
	flag.StringVar(&config.WebPath, "webpath", "/", "Specify added url http://host:port/<path> Default: /")
	flag.StringVar(&config.UpdateKey, "updatekey", "changeme", "Specify key to API update")
	flag.StringVar(&config.DashName, "dashname", "", "Dashboard name")
	flag.StringVar(&configfile, "configfile", "dashgoat.yaml", "Name of configfile")
	flag.StringVar(&buddy_cli.Url, "buddyurl", "", "Buddy url")
	flag.StringVar(&buddy_cli.Key, "buddykey", "", "Buddy update key, empty for same key")
	flag.StringVar(&buddy_cli.Name, "buddyname", "", "Buddy name")
	flag.StringVar(&nsconfig, "nsconfig", "", "Configure buddies via DNS/k8s namespace")
	flag.Parse()

	config.ReadEnv()
	err := config.InitConfig(configfile)
	if err != nil {
		fmt.Println(err)
	}

	pathStartsWith := strings.HasPrefix(config.WebPath, "/")
	if !pathStartsWith {
		config.WebPath = "/" + config.WebPath
	}

	useOS := false
	// Serving embedded static files
	// if len(os.Args) > 1 && os.Args[1] == "live" {
	// 	useOS = true
	// }
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))

	e.HideBanner = true

	if config.WebLog == "on" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	if !config.DisableMetrics {
		e.Use(echoprometheus.NewMiddleware("dashgoat"))
		serviceStateCollector = NewServiceStateCollector()
		prometheus.MustRegister(serviceStateCollector)
		e.GET(add2url(config.WebPath, "/metrics"), echoprometheus.NewHandler())
	}

	e.GET(add2url(config.WebPath, "/metricshistory/:serviceid/:hours"), getMetricsHistory)

	e.POST(add2url(config.WebPath, "/update"), updateStatus)
	e.GET(add2url(config.WebPath, "/status/:id"), getStatus)
	e.GET(add2url(config.WebPath, "/status/list"), getStatusList)
	e.GET(add2url(config.WebPath, "/status/listmso"), getStatusListMSO)
	e.GET(add2url(config.WebPath, "/list/:serviceitem"), getUniq)
	e.GET(add2url(config.WebPath, "/servicefilter/:item/:itemval"), serviceFilter)
	e.DELETE(add2url(config.WebPath, "/service/:id"), deleteServiceHandler)
	e.GET(add2url(config.WebPath, "/health"), health)

	printWelcome()

	go lostProbeTimer()
	go ttlHousekeeping()
	go findBuddy(config.BuddyHosts)

	// Start server
	e.Logger.Fatal(e.Start(config.IPport))

}

func printWelcome() {
	fmt.Println("Starting dashGoat " + readHostFacts().DashGoatVersion)
	fmt.Println("Dashboard name: " + readHostFacts().DashName + " ")
	fmt.Println("Go: " + readHostFacts().GoVersion + " ")
	fmt.Println("Labstack Echo: " + echo.Version + " ")
}
