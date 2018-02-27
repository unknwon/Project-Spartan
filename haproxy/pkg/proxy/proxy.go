// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	log "gopkg.in/clog.v1"
)

// Proxy is the core part of HA proxy, it maintains the list of server end points
// and has the right to pick active server.
type Proxy struct {
	// List of available server end points
	endPoints []string
	// Time interval of two health check rounds
	healthCheckInterval time.Duration
	// Timeout duration of a health check request
	healthCheckTimeout time.Duration
	// HTTP client for health check
	healthCheckClient *http.Client

	// RW mutex for active server reverse proxy object
	proxyLocker sync.RWMutex
	// Active server end point
	activeEndPoint string
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
			endPoints:           endPoints,
			healthCheckInterval: healthCheckInterval,
			healthCheckTimeout:  healthCheckTimeout,
			healthCheckClient: &http.Client{
				Timeout: healthCheckTimeout,
			},
			activeEndPoint: endPoints[0],
		}

		proxyInstance.HealthCheck()
		proxyInstance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   proxyInstance.activeEndPoint,
		})
	})
	return proxyInstance
}

var healthCheckCount int64 = 1

func (p *Proxy) sendHealthCheckRequest(endPoint string) bool {
	resp, err := p.healthCheckClient.Get("http://" + endPoint + "/healthcheck")
	if err != nil {
		log.Error(2, "Fail to perform health check for '%s': %v", endPoint, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(2, "Fail to perform health check for '%s': status code is %d not 200", endPoint, resp.StatusCode)
		return false
	}
	return true
}

// HealthCheck sends out health check requests to all server end points
func (p *Proxy) HealthCheck() {
	log.Trace("[%d] Health check started...", healthCheckCount)

	for _, endPoint := range p.endPoints {
		if p.sendHealthCheckRequest(endPoint) {
			// No need to recreate same reverse proxy object if active end point is already it
			if endPoint == p.activeEndPoint {
				log.Trace("[HC] Active server '%s' still healthy", endPoint)
				break
			}

			p.proxyLocker.Lock()
			p.activeEndPoint = endPoint
			proxyInstance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   proxyInstance.activeEndPoint,
			})
			p.proxyLocker.Unlock()
			log.Trace("[HC] Active server changed to: %s", endPoint)
			break
		}
	}

	time.AfterFunc(p.healthCheckInterval, func() { p.HealthCheck() })
	log.Trace("[%d] Health check finished.", healthCheckCount)
	healthCheckCount++
}
