// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"net/http"
	"time"

	log "gopkg.in/clog.v1"

	"github.com/Unknwon/Project-Spartan/haproxy/pkg/registry"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/setting"
)

var haproxyRegistry, serverRegistry *registry.Registry

func init() {
	healthCheckClient = &http.Client{
		Timeout: setting.HealthCheck.Timeout,
	}
	haproxyRegistry = registry.NewRegistry(setting.HAProxy.EndPoints)
	serverRegistry = registry.NewRegistry(setting.Server.EndPoints)

	go HealthCheck()
}

var healthCheckClient *http.Client
var healthCheckCount int64 = 1

func sendHealthCheckRequest(server *registry.Instance) bool {
	resp, err := healthCheckClient.Get("http://" + server.Address + "/healthcheck")
	if err != nil {
		if _, ok := err.(net.Error); ok {
			log.Warn("[HC] Instance '%s' is down", server)
		} else {
			log.Error(2, "Fail to perform health check for '%s': %v", server, err)
		}
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(2, "Fail to perform health check for '%s': status code is %d not 200", server, resp.StatusCode)
		return false
	}
	return true
}

func HealthCheck() {
	log.Trace("[%d] Health check started...", healthCheckCount)

	for _, servers := range [][]*registry.Instance{
		haproxyRegistry.Servers,
		serverRegistry.Servers,
	} {
		for _, server := range servers {
			if sendHealthCheckRequest(server) {
				server.Status = registry.STATUS_RUNNING
			} else {
				server.Status = registry.STATUS_DOWN
			}
		}
	}

	time.AfterFunc(setting.HealthCheck.Interval, HealthCheck)
	log.Trace("[%d] Health check finished.", healthCheckCount)
	healthCheckCount++
}
