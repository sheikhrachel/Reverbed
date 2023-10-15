# Reverbed

## Description

This is my lil Go http server for the reverbed project, built with workbench.

## Repo Structure

- `api_common`: a general purpose package with all of my util code for http routers, aws, caching, etc.
- `handlers`: a package for all the http handlers and the route registry
- `model`: a package for all the data models used throughout the service
- `*-service`: an isolated package of business logic functions that can be deployed independently, or as a part of a larger service like `handlers`
- `swaggerui`: a simple swagger ui page that can simplify documentation based on yaml or json openapi files

## Running Locally

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)

### Commands

- `make`: run the server locally.  spins up mock services outlined in `docker-compose.yaml`
  - This starts the server on port 8080
  - Additional AWS infra can be connected to in `api_common/aws_utils/aws_client.go`
