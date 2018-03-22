// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

var (
	port      = flag.Int("port", 8002, "Listening port number for reseller portal")
	name      = flag.String("name", os.Getenv("SPARTAN-RPORTAL-NAME"), "Code name for reseller portal instance")
	mysqlHost = flag.String("mysql-host", os.Getenv("SPARTAN-MYSQL-HOST"), "Code name for reseller portal instance")
)

func init() {
	flag.Parse()
	log.New(log.CONSOLE, log.ConsoleConfig{})
}

func main() {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("ui/dist"))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "ui/dist",
	}))
	m.Use(func(c *macaron.Context) {
		c.Data["CodeName"] = *name
		c.Resp.Header().Add("X-Spartan-Server", *name)
	})

	m.Get("/", Home)
	m.Get("/healthcheck", HealthCheck)
	m.Get("/metadata", MetaData)

	m.Group("/api", func() {
		m.Combo("/items").Get(ListItems).Post(AddItem)
		m.Delete("/items/:id", DeleteItem)
	})

	log.Info("Instance: %s, running on port %d", *name, *port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), m); err != nil {
		log.Fatal(2, "Fail to start server: %v", err)
	}
}
