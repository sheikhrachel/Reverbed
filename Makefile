CURR_DIR            := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
ROOT                := $(shell pwd)
SERVICE_NAME        := reverbed
PKG                 := ${ROOT}
CMD                 := ${ROOT}/main.go
VERSION             := 1.21
DOCKER_SRC          := public.ecr.aws/docker/library/golang
GO_TEST_OUTFILE     := ./c.out
GO_HTML_COV         := ./coverage.html

.PHONY: default
default: compose

###############################
# make || make compose
# will run this service and
# enable you to hit endpoints
# on localhost:4000
#
# We override 4000 to 8080
# in the docker-compose.yaml
###############################
compose:
ifndef fast
	docker-compose -f docker-compose.yaml rm -f
	docker-compose -f docker-compose.yaml build --no-cache
endif
	docker-compose -f docker-compose.yaml up
	$(info "Service is UP!")
	$(info "=> curl http://localhost:4000/health")
shutdown:
	docker-compose -f docker-compose.yaml down
