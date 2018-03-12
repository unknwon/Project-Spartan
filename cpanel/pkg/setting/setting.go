// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"strings"
	"time"

	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
)

var (
	HAProxy struct {
		EndPoints []string
	}
	Server struct {
		EndPoints []string
	}
	Database struct {
		EndPoints []string
	}
	HealthCheck struct {
		Interval time.Duration
		Timeout  time.Duration
	}

	Config *ini.File
)

func init() {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	var err error
	Config, err = ini.Load("conf/app.ini", "conf/custom.ini")
	if err != nil {
		log.Fatal(2, "Fail to load configuration: %v", err)
	}
	Config.NameMapper = ini.AllCapsUnderscore

	if err = Config.Section("haproxy").MapTo(&HAProxy); err != nil {
		log.Fatal(2, "Fail to map HAProxy settings: %v", err)
	}
	log.Info("HAProxy end points: %s", strings.Join(HAProxy.EndPoints, ", "))

	if err = Config.Section("server").MapTo(&Server); err != nil {
		log.Fatal(2, "Fail to map Server settings: %v", err)
	}
	log.Info("Server end points: %s", strings.Join(Server.EndPoints, ", "))

	if err = Config.Section("database").MapTo(&Database); err != nil {
		log.Fatal(2, "Fail to map Database settings: %v", err)
	}
	log.Info("Database end points: %s", strings.Join(Database.EndPoints, ", "))

	if err = Config.Section("health_check").MapTo(&HealthCheck); err != nil {
		log.Fatal(2, "Fail to map HealthCheck settings: %v", err)
	}
	log.Info("Health check (interval/timeout): %s/%s", HealthCheck.Interval, HealthCheck.Timeout)
}
