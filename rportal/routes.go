// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"

	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

func HealthCheck(c *macaron.Context) {
	c.JSON(200, map[string]interface{}{
		"Status": "OK",
	})
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
