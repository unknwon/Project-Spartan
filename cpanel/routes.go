// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/docker"
)

func Home(c *macaron.Context) {
	c.HTML(200, "app")
}

func Dashboard(c *macaron.Context) {
	c.JSON(200, map[string]interface{}{
		"haproxies": haproxyRegistry.Instances,
		"servers":   serverRegistry.Instances,
		"databases": databaseRegistry.Instances,
	})
}

func StartServer(c *macaron.Context) {
	// Note: Currently only support application runs on Docker container.
	in, err := serverRegistry.InstanceByName(c.Query("name"))
	if err != nil {
		c.PlainText(422, []byte(err.Error()))
		return
	}

	switch {
	case strings.Contains(in.Name, "docker"):
		if err := docker.CreateContainer(in.Name, in.Address, "rportal:latest"); err != nil {
			c.PlainText(500, []byte(err.Error()))
			return
		}
	default:
		c.PlainText(422, []byte("Only application runs on Docker is supported"))
		return
	}
	c.Status(204)
}

func ShutdownServer(c *macaron.Context) {
	// Note: Currently only support application runs on Docker container.
	in, err := serverRegistry.InstanceByName(c.Query("name"))
	if err != nil {
		c.PlainText(422, []byte(err.Error()))
		return
	}

	switch {
	case strings.Contains(in.Name, "docker"):
		if err := docker.ShutdownContainer(in.Name); err != nil {
			c.PlainText(500, []byte(err.Error()))
			return
		}
	default:
		c.PlainText(422, []byte("Only application runs on Docker is supported"))
		return
	}
	c.Status(204)
}
