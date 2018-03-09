// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("ui/dist"))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "ui/dist",
	}))

	m.Get("/", Home)

	m.Group("/api", func() {
		m.Get("/dashboard", Dashboard)
	})

	log.Info("Instance: %s, running on port %d", "cpanel", 7777)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 7777), m); err != nil {
		log.Fatal(2, "Fail to start server: %v", err)
	}
}
