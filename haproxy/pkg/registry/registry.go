// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package registry

import (
	"strings"
)

type Status string

const (
	STATUS_UNKNOWN Status = "unknown"
	STATUS_RUNNING Status = "running"
	STATUS_DOWN    Status = "down"
)

// Instance contains information of an instance's code name and address.
type Instance struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  Status `json:"status"`
}

func (s *Instance) String() string {
	return s.Name + "/" + s.Address
}

// Registry maintains a list of servers.
type Registry struct {
	Servers []*Instance
}

// NewRegistry parses raw input of server metadata and returns a registry maintains the list.
// The raw input format should be:
//		["rportal-local-1/localhost:8002", "rportal-docker-1/localhost:9002"]
func NewRegistry(inputs []string) *Registry {
	r := &Registry{
		Servers: make([]*Instance, len(inputs)),
	}
	for i, input := range inputs {
		fields := strings.Split(input, "/")
		r.Servers[i] = &Instance{
			Name:    fields[0],
			Address: fields[1],
			Status:  STATUS_UNKNOWN,
		}
	}
	return r
}
