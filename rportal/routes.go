// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"gopkg.in/macaron.v1"
)

func HealthCheck(c *macaron.Context) {
	c.JSON(200, map[string]interface{}{
		"Status": "OK",
	})
}

func Home(c *macaron.Context) {
	c.PlainText(200, []byte("Serving content from "+*name))
}
