// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Note: This package is not thread-safe!
package registry

import (
	"fmt"
	"strings"
	"sync"
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
	Score   int    `json:"score"`
}

func (s *Instance) String() string {
	return s.Name + "/" + s.Address
}

// Registry maintains a list of servers.
type Registry struct {
	locker    sync.RWMutex
	Instances []*Instance
}

// NewRegistry parses raw input of instance metadata and returns a registry maintains the list.
// The raw input format should be:
//		["rportal-local-1/localhost:8002", "rportal-docker-1/localhost:9002"]
func NewRegistry(inputs []string) *Registry {
	r := &Registry{
		Instances: make([]*Instance, len(inputs)),
	}
	for i, input := range inputs {
		fields := strings.SplitN(input, "/", 2)
		r.Instances[i] = &Instance{
			Name:    fields[0],
			Address: fields[1],
			Status:  STATUS_UNKNOWN,
		}
	}
	return r
}

// InstanceByName returns an instance object by given name.
// It returns an error if no instance found associated with the name.
func (r *Registry) InstanceByName(name string) (*Instance, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	for i := range r.Instances {
		if r.Instances[i].Name == name {
			return r.Instances[i], nil
		}
	}
	return nil, fmt.Errorf("instance '%s' not found", name)
}

func (r *Registry) SetInstanceAddress(name, address string) error {
	in, err := r.InstanceByName(name)
	if err != nil {
		return err
	}

	r.locker.Lock()
	defer r.locker.Unlock()

	in.Address = address
	return nil
}
