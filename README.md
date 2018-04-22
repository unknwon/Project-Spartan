# Project Spartan

The Project Spartan is a hybrid-infrastructure and research-oriented HA (High Availability) prototype for IT 485.

## Project Status

The goal of this project is more like PoC (Proof of Concept) rather than comprehensive implementation.

## Prerequisites

All prerequisites are for local deployment (i.e. on your Mac) only, cloud deployment is much easier that only needs compiled binaries of all components (i.e. haproxy, rportal, and cpanel).

_Built based on such version, other versions are not tested_

- Go 1.10 
- Docker for Mac 18.03.0-ce
- AWS CLI 1.15.0
- Google Cloud SDK 196.0.0
- VirtualBox 5.2.8
- Vagrant 2.0.3
- CoreDNS 1.1.1
- MySQL 5.7.22

## System Overview

![Network Topology](docs/images/Network%20Topology.jpg)

### Components

- CoreDNS: Response DNS query to avoid SPOF of haproxy
- haproxy (e.g. `haproxy-vagrant-1`, `haproxy-vagrant-2`): 
    - Load balance based on response time, CPU load, and memory usage
    - Request health checks to rportal
    - Response health check from cpanel
- rportal (e.g. `rportal-local-1`, `rportal-docker-1`, `rportal-aws-us-east-1`, `rportal-gcp-us-east-1`):
    - Serve web pages
    - Response to health checks from haproxy and cpanel
    - Process and store data to MySQL
- cpanel: Monitor and manipulate other components
- MySQL (e.g. `mysql-aws-us-east-1`): Store data

## Disclaimers

The security is not part of main focus of the system, only basic level of security is considered.

## License

This project is under MIT License **with exception of any academic honesty violation**. See the [LICENSE](LICENSE) file for the full license text.