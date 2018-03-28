// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"encoding/json"
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

type HealthCheckResponse struct {
	Status      string
	CPULoad     int
	MemoryUsage int
}

// sendHealthCheckRequest sends the actual health check request and calculates the health score.
// It also returns an additional boolean to indicate if the request was succeed.
func (p *Proxy) sendHealthCheckRequest(in *registry.Instance) (int, bool) {
	start := time.Now()

	resp, err := p.healthCheckClient.Get("http://" + in.Address + "/healthcheck")
	if err != nil {
		if _, ok := err.(net.Error); ok {
			log.Warn("[HC] Instance '%s' is down", in)
		} else {
			log.Error(2, "Fail to perform health check for '%s': %v", in, err)
		}
		return -1, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(2, "Fail to perform health check for '%s': status code is %d not 200", in, resp.StatusCode)
		return -1, false
	}

	elapsed := int(time.Since(start).Nanoseconds() / 1000000)

	// Parse response and calculate health score
	var hcResp HealthCheckResponse
	if err = json.NewDecoder(resp.Body).Decode(&hcResp); err != nil {
		log.Error(2, "Fail to decode health check response for '%s': %v", in, err)
		return -1, false
	}

	score := (elapsed + hcResp.CPULoad + hcResp.MemoryUsage) / 10
	log.Trace("[HC] Instance '%s' score: %d (%dms/%d/%d)", in, score, elapsed, hcResp.CPULoad, hcResp.MemoryUsage)

	return score, true
}

// HealthCheck sends out health check requests to all server end points.
func (p *Proxy) HealthCheck() {
	log.Trace("[%d] Health check started...", healthCheckCount)
	defer func() {
		log.Trace("[%d] Health check finished.", healthCheckCount)
		healthCheckCount++
		time.AfterFunc(p.healthCheckInterval, func() { p.HealthCheck() })
	}()

	var candidateServer *registry.Instance
	for _, in := range p.registry.Instances {
		score, ok := p.sendHealthCheckRequest(in)
		in.Score = score

		if !ok {
			in.Status = registry.STATUS_DOWN
			continue
		}

		if candidateServer == nil ||
			candidateServer.Score > in.Score {
			candidateServer = in
		}
	}

	// In case no instance is healthy
	if candidateServer == nil {
		log.Error(2, "ALERT: All instances are down!!!")
		return
	}

	// No need to recreate same reverse proxy object if active end point is already it
	if candidateServer.Name == p.activeServer.Name {
		log.Trace("[HC] Instance '%s' is still the active server", p.activeServer)
		return
	}

	p.proxyLocker.Lock()
	p.activeServer.Name = candidateServer.Name
	p.activeServer.Address = candidateServer.Address
	proxyInstance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   proxyInstance.activeServer.Address,
	})
	p.proxyLocker.Unlock()
	log.Info("[HC] Active server changed to: %s", p.activeServer)
}
