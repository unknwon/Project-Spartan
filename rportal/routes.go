// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"math/rand"
	"time"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

var (
	responseDelay int // Between 0-500
	cpuLoad       int // Between 0-100
	memoryUsage   int // Between 0-100
)

// Randomize response time, CPU load, memory usage at start.
func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	responseDelay = r1.Intn(500)
	cpuLoad = rand.Intn(100)
	memoryUsage = r1.Intn(100)
}

func HealthCheck(c *macaron.Context) {
	time.Sleep(time.Duration(responseDelay) * time.Millisecond)

	c.JSON(200, map[string]interface{}{
		"Status":      "OK",
		"CPULoad":     cpuLoad,
		"MemoryUsage": memoryUsage,
	})
}

func MetaData(c *macaron.Context) {
	responseDelay = c.QueryInt("responseDelay")
	cpuLoad = c.QueryInt("cpuLoad")
	memoryUsage = c.QueryInt("memoryUsage")
}

func Home(c *macaron.Context) {
	c.HTML(200, "app")
}

func ListItems(c *macaron.Context) {
	items := make([]*Reseller, 0, 5)
	if err := x.Find(&items).Error; err != nil {
		log.Error(2, "Fail to read items: %v", err)
		c.Status(500)
		return
	}

	c.JSON(200, items)
}

func AddItem(c *macaron.Context) {
	r := &Reseller{}
	data, err := c.Req.Body().Bytes()
	if err != nil {
		log.Error(2, "Fail to read request body: %v", err)
		c.Status(500)
		return
	}
	if err = json.Unmarshal(data, r); err != nil {
		log.Error(2, "Fail to parse request body: %v", err)
		c.Status(500)
		return
	}
	if err := x.Create(r).Error; err != nil {
		log.Error(2, "Fail to add item: %v", err)
		c.Status(500)
		return
	}

	c.JSON(201, r)
}

func DeleteItem(c *macaron.Context) {
	if err := x.Delete(new(Reseller), "id = ?", c.ParamsInt64("id")).Error; err != nil {
		log.Error(2, "Fail to delete item: %v", err)
		c.Status(500)
		return
	}
	c.Status(204)
}
