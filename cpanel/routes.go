// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/awsec2"
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
	case strings.Contains(in.Name, "-docker-"):
		err = docker.CreateContainer(in.Name, in.Address, "rportal:latest")
	case strings.Contains(in.Name, "-aws-"):
		err = awsec2.StartInstance(in.Name)

		time.AfterFunc(10*time.Second, func() {
			for {
				ip, err := awsec2.GetInstancePublicIPv4(in.Name)
				if err == nil {
					if ip != "None" {
						serverRegistry.SetInstanceAddress(in.Name, ip+":8002")
						break
					}
				}

				time.Sleep(1 * time.Second)
			}
		})

	default:
		c.PlainText(422, []byte("Application runs on given infrastructure is not supported"))
		return
	}

	if err != nil {
		c.PlainText(500, []byte(err.Error()))
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
	case strings.Contains(in.Name, "-docker-"):
		err = docker.ShutdownContainer(in.Name)
	case strings.Contains(in.Name, "-aws-"):
		err = awsec2.ShutdownInstance(in.Name)

	default:
		c.PlainText(422, []byte("Application runs on given infrastructure is not supported"))
		return
	}

	if err != nil {
		c.PlainText(500, []byte(err.Error()))
		return
	}
	c.Status(204)
}
