// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "gopkg.in/clog.v1"
)

var x *gorm.DB
var xs *gorm.DB

type Reseller struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Person  string `json:"person"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func init() {
	var err error
	x, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/rportal", *mysqlUser, *mysqlPassword, *mysqlHost))
	if err != nil {
		log.Fatal(2, "Fail to open database connection: %v", err)
	}

	xs, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/rportal", *mysqlUser, *mysqlPassword, *mysqlSlaveHost))
	if err != nil {
		log.Fatal(2, "Fail to open slave database connection: %v", err)
	}

	if x.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&Reseller{}).Error != nil {
		log.Fatal(2, "Fail to auto migrate database tables: %v", err)
	}
}

// getDB returns slave DB object when master fails.
func getDB() *gorm.DB {
	e := x
	x.HasTable("resellers")
	if err := x.DB().Ping(); err != nil {
		log.Warn("Master database is currently unavailable: %v", err)
		e = xs
	}
	return e
}
