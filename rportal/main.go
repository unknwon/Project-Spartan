// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

var (
	port = flag.Int("port", 8002, "Listening port number for reseller portal")
	name = flag.String("name", "undefined", "Code name for reseller portal instance")
)

func init() {
	flag.Parse()
	log.New(log.CONSOLE, log.ConsoleConfig{})
}

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	m.Get("/", Home)
	m.Get("/healthcheck", HealthCheck)

	log.Info("Instance: %s, running on port %d", *name, *port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), m); err != nil {
		log.Fatal(2, "Fail to start server: %v", err)
	}
}
