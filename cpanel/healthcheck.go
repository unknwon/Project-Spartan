// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "gopkg.in/clog.v1"

	"github.com/Unknwon/Project-Spartan/haproxy/pkg/registry"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/setting"
)

var haproxyRegistry, serverRegistry, databaseRegistry *registry.Registry
var dbCoons = map[string]*gorm.DB{}

func init() {
	healthCheckClient = &http.Client{
		Timeout: setting.HealthCheck.Timeout,
	}
	haproxyRegistry = registry.NewRegistry(setting.HAProxy.EndPoints)
	serverRegistry = registry.NewRegistry(setting.Server.EndPoints)
	databaseRegistry = registry.NewRegistry(setting.Database.EndPoints)

	// Setup database connections
	for _, in := range databaseRegistry.Instances {
		x, err := gorm.Open("mysql", fmt.Sprintf("root:@tcp(%s)/rportal", in.Address))
		if err != nil {
			log.Fatal(2, "Fail to open database connection: %v", err)
		}
		dbCoons[in.Name] = x
	}

	go HealthCheck()
}

var healthCheckClient *http.Client
var healthCheckCount int64 = 1

func sendHealthCheckRequest(in *registry.Instance) bool {
	resp, err := healthCheckClient.Get("http://" + in.Address + "/healthcheck")
	if err != nil {
		if _, ok := err.(net.Error); ok {
			log.Warn("[HC] Instance '%s' is down", in)
		} else {
			log.Error(2, "Fail to perform health check for '%s': %v", in, err)
		}
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(2, "Fail to perform health check for '%s': status code is %d not 200", in, resp.StatusCode)
		return false
	}
	return true
}

func HealthCheck() {
	log.Trace("[%d] Health check started...", healthCheckCount)

	for _, ins := range [][]*registry.Instance{
		haproxyRegistry.Instances,
		serverRegistry.Instances,
	} {
		for _, in := range ins {
			if sendHealthCheckRequest(in) {
				in.Status = registry.STATUS_RUNNING
			} else {
				in.Status = registry.STATUS_DOWN
			}
		}
	}

	for _, in := range databaseRegistry.Instances {
		if err := dbCoons[in.Name].DB().Ping(); err != nil {
			log.Error(2, "Fail to perform health check for '%s': %v", in, err)
			in.Status = registry.STATUS_DOWN
		} else {
			in.Status = registry.STATUS_RUNNING
		}
	}

	time.AfterFunc(setting.HealthCheck.Interval, HealthCheck)
	log.Trace("[%d] Health check finished.", healthCheckCount)
	healthCheckCount++
}
