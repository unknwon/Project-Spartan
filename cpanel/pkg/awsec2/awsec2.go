// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package awsec2

import (
	"fmt"
	"strings"

	"github.com/Unknwon/Project-Spartan/cpanel/pkg/command"
)

// GetInstanceID returns the instance ID by given key name.
func GetInstanceID(kname string) (string, error) {
	stdout, err := command.New("aws", "ec2", "describe-instances",
		"--filter", "Name=key-name,Values=rportal-aws-us-east-1",
		"--query", "Reservations[].Instances[].[InstanceId]",
		"--output=text").Run()
	return strings.Fields(stdout)[0], err
}

// StartInstance starts an instance with given key name.
func StartInstance(kname string) error {
	id, err := GetInstanceID(kname)
	if err != nil {
		return fmt.Errorf("GetInstanceID: %v", err)
	}

	_, err = command.New("aws", "ec2", "start-instances",
		"--instance-id", id).Run()
	return err
}

// GetInstancePublicIPv4 returns the public IPv4 address of the instance by given key name.
func GetInstancePublicIPv4(kname string) (string, error) {
	stdout, err := command.New("aws", "ec2", "describe-instances",
		"--filter", "Name=key-name,Values=rportal-aws-us-east-1",
		"--query", "Reservations[].Instances[].[PublicIpAddress]",
		"--output=text").Run()
	return strings.Fields(stdout)[0], err
}

// ShutdownInstance stops an instance with given key name.
func ShutdownInstance(kname string) error {
	id, err := GetInstanceID(kname)
	if err != nil {
		return fmt.Errorf("GetInstanceID: %v", err)
	}

	_, err = command.New("aws", "ec2", "stop-instances",
		"--instance-id", id).Run()
	return err
}
