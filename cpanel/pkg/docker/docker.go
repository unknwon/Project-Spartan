// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package docker

import (
	"strings"
)

// CreateContainer creates a new container with given name based on given image name.
// It tries to remove the existing container with same name before creating a new one.
func CreateContainer(cname, caddr, iname string) error {
	NewCommand("rm", cname).Run()

	_, err := NewCommand("run", "-d",
		"-e", "SPARTAN-RPORTAL-NAME="+cname,
		"-e", "SPARTAN-MYSQL-HOST=docker.for.mac.host.internal",
		"-p", strings.Split(caddr, ":")[1]+":8002",
		"--name", cname,
		iname).Run()
	return err
}

// ShutdownContainer stops a container with given name.
func ShutdownContainer(cname string) error {
	_, err := NewCommand("stop", cname).Run()
	return err
}
