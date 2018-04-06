// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package docker

import (
	"strings"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/command"
)

// CreateContainer creates a new container with given name based on given image name.
// It tries to remove the existing container with same name before creating a new one.
func CreateContainer(cname, caddr, iname string) error {
	command.New("docker", "rm", cname).Run()

	_, err := command.New("docker", "run", "-d",
		"-e", "SPARTAN-RPORTAL-NAME="+cname,
		"-e", "SPARTAN-MYSQL-HOST=docker.for.mac.host.internal",
		"-p", strings.Split(caddr, ":")[1]+":8002",
		"--name", cname,
		iname).Run()
	return err
}

// ShutdownContainer stops a container with given name.
func ShutdownContainer(cname string) error {
	_, err := command.New("docker", "stop", cname).Run()
	return err
}
