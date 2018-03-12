// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	log "gopkg.in/clog.v1"

	"github.com/Unknwon/Project-Spartan/haproxy/pkg/registry"
)

// Proxy is the core part of HA proxy, it maintains the list of server end points
// and has the right to pick active server.
type Proxy struct {
	registry *registry.Registry
	// Time interval of two health check rounds
	healthCheckInterval time.Duration
	// Timeout duration of a health check request
	healthCheckTimeout time.Duration
	// HTTP client for health check
	healthCheckClient *http.Client

	// RW mutex for active server reverse proxy object
	proxyLocker sync.RWMutex
	// Active server end point
	activeServer *registry.Instance
	// Reverse proxy object corresponding to the active server end point
	proxy *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxyLocker.RLock()
	defer p.proxyLocker.RUnlock()

	p.proxy.ServeHTTP(w, r)
}

var (
	proxyInstance *Proxy
	once          sync.Once
)

// NewProxy creates and only creates single instance of Proxy object with given list
// of server end points, health check interval and timeout.
func NewProxy(endPoints []string, healthCheckInterval, healthCheckTimeout time.Duration) *Proxy {
	if len(endPoints) == 0 {
		panic("expect at least one end points, but got zero")
	}

	once.Do(func() {
		proxyInstance = &Proxy{
			registry:            registry.NewRegistry(endPoints),
			healthCheckInterval: healthCheckInterval,
			healthCheckTimeout:  healthCheckTimeout,
			healthCheckClient: &http.Client{
				Timeout: healthCheckTimeout,
			},
		}
		proxyInstance.activeServer = &registry.Instance{
			Name:    proxyInstance.registry.Instances[0].Name,
			Address: proxyInstance.registry.Instances[0].Address,
		}

		proxyInstance.HealthCheck()
		proxyInstance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   proxyInstance.activeServer.Address,
		})
	})
	return proxyInstance
}

var healthCheckCount int64 = 1

func (p *Proxy) sendHealthCheckRequest(server *registry.Instance) bool {
	resp, err := p.healthCheckClient.Get("http://" + server.Address + "/healthcheck")
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

// HealthCheck sends out health check requests to all server end points
func (p *Proxy) HealthCheck() {
	log.Trace("[%d] Health check started...", healthCheckCount)

	for _, in := range p.registry.Instances {
		if p.sendHealthCheckRequest(in) {
			// No need to recreate same reverse proxy object if active end point is already it
			if in.Name == p.activeServer.Name {
				log.Trace("[HC] Active server '%s' still healthy", in)
				break
			}

			p.proxyLocker.Lock()
			p.activeServer.Name = in.Name
			p.activeServer.Address = in.Address
			proxyInstance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   proxyInstance.activeServer.Address,
			})
			p.proxyLocker.Unlock()
			log.Info("[HC] Active server changed to: %s", in)
			break
		}
	}

	time.AfterFunc(p.healthCheckInterval, func() { p.HealthCheck() })
	log.Trace("[%d] Health check finished.", healthCheckCount)
	healthCheckCount++
}
