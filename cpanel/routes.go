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
	"github.com/Unknwon/Project-Spartan/cpanel/pkg/gcpvm"
	"github.com/Unknwon/Project-Spartan/cpanel/pkg/setting"
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

			// Update configuration file
			setting.Config.Section("server").Key("END_POINTS").SetValue(strings.Join(serverRegistry.List(), ", "))
			setting.Config.SaveTo(setting.CUSTOM_CONF_PATH)
		})
	case strings.Contains(in.Name, "-gcp-"):
		err = gcpvm.StartInstance(in.Name)
		if err != nil {
			break
		}

		var ip string
		ip, err = gcpvm.GetInstancePublicIPv4(in.Name)
		if err == nil {
			serverRegistry.SetInstanceAddress(in.Name, ip+":8002")

			// Update configuration file
			setting.Config.Section("server").Key("END_POINTS").SetValue(strings.Join(serverRegistry.List(), ", "))
			setting.Config.SaveTo(setting.CUSTOM_CONF_PATH)
		}

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
	case strings.Contains(in.Name, "-gcp-"):
		err = gcpvm.ShutdownInstance(in.Name)

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
