// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gcpvm

import (
	"github.com/tidwall/gjson"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/command"
)

// StartInstance starts an instance with given name.
func StartInstance(name string) error {
	_, err := command.New("gcloud", "compute", "instances", "start", name).Run()
	return err
}

// GetInstancePublicIPv4 returns the public IPv4 address of the instance by given name.
func GetInstancePublicIPv4(name string) (string, error) {
	stdout, err := command.New("gcloud", "compute", "instances", "list",
		"--filter", "name=('"+name+"')",
		"--format=json").Run()

	return gjson.Get(stdout, "0.networkInterfaces.0.accessConfigs.0.natIP").String(), err
}

// ShutdownInstance stops an instance with given name.
func ShutdownInstance(name string) error {
	_, err := command.New("gcloud", "compute", "instances", "stop", name).Run()
	return err
}
