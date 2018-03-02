// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package registry

import (
	"strings"
)

// Server contains information of a server's code name and address.
type Server struct {
	Name    string
	Address string
}

func (s *Server) String() string {
	return s.Name + "/" + s.Address
}

// Registry maintains a list of servers.
type Registry struct {
	Servers []*Server
}

// NewRegistry parses raw input of server metadata and returns a registry maintains the list.
// The raw input format should be:
//		["rportal-local-1/localhost:8002", "rportal-docker-1/localhost:9002"]
func NewRegistry(inputs []string) *Registry {
	r := &Registry{
		Servers: make([]*Server, len(inputs)),
	}
	for i, input := range inputs {
		fields := strings.Split(input, "/")
		r.Servers[i] = &Server{
			Name:    fields[0],
			Address: fields[1],
		}
	}
	return r
}
