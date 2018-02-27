// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	log "gopkg.in/clog.v1"

	"github.com/Unknwon/Project-Spartan/haproxy/pkg/proxy"
	"github.com/Unknwon/Project-Spartan/haproxy/pkg/setting"
)

var (
	port = flag.Int("port", 8000, "Listening port number for HA proxy")
	name = flag.String("name", "undefined", "Code name for HA proxy instance")
)

func init() {
	flag.Parse()
}

func main() {
	log.Info("Instance: %s, running on port %d", *name, *port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), &proxyHandler{
		Proxy: proxy.NewProxy(setting.Server.EndPoints, setting.HealthCheck.Interval, setting.HealthCheck.Timeout),
	}); err != nil {
		log.Fatal(2, "Fail to start proxy: %v", err)
	}
}

type proxyHandler struct {
	*proxy.Proxy
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Spartan-Proxy", *name)
	h.Proxy.ServeHTTP(w, r)
}
