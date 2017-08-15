
.PHONY: all build-all build-app build-image repo-warning clean-app push clean golang-build-container echovars vet clean-image clean-all very-clean
#
# This is the Makefile to build the golang components.
# It is normally run in a docker container for golang builds.
#
# 3/14/2016 mln
#
# golog compile via docker container:
# https://hub.docker.com/_/golang/
#
# Vars for version and build
VERSION ?= 1.0.0
#
# do some repo evals... want everyone to use their own, unless building for the team
#
DEFAULT_REPO := "samsung-bogus"
DOCKER_REPO ?= DEFAULT_REPO
REPO := $(DOCKER_REPO)
#
#
DOCKER_IMAGE := jsonsvalidator-img
SOURCE_APP := jsonsvalidator-linux

repo-warning:
	@echo "============================================================="
	@echo "make setup:$(MAKE):$(MAKECMDGOALS):$(MAKEFLAGS):"
	@if  [ $(DOCKER_REPO) =  $(DEFAULT_REPO) ]; then \
        	echo "+++++++++++++++++++++++++++++++++++++++++++++++++++"; \
        	echo "  You have not changed DOCKER_REPO from: $(DOCKER_REPO)"; \
        	echo "  You MUST set DOCKER_REPO in your environment"; \
        	echo "  or directly in this Makefile unless you are"; \
        	echo "  building for the group"; \
        	echo "+++++++++++++++++++++++++++++++++++++++++++++++++++"; \
        	false; \
    	else \
        	echo "+++++++++++++++++++++++++++++++++++++++++++++++++++"; \
        	echo "  Your DOCKER_REPO is set to: $(DOCKER_REPO)"; \
        	echo "  Please execute 'make all' to build"; \
        	echo "+++++++++++++++++++++++++++++++++++++++++++++++++++"; \
    	fi

all: build-all

build-all: build-image
#
# DOCKER IMAGE BUILD
#
build-image: $(DOCKER_IMAGE)

$(DOCKER_IMAGE): Dockerfile $(SOURCE_APP) ca-certificates.crt
	@echo "building $(DOCKER_IMAGE) $(VERSION)"
	docker build -t $(REPO)/$(DOCKER_IMAGE):$(VERSION) --rm=true --force-rm=true .
	@touch $@
	@docker images $(REPO)/$(DOCKER_IMAGE)


push:
	docker push $(REPO)/$(DOCKER_IMAGE):$(VERSION)

clean: clean-all

clean-all: clean-image

clean-image:
	-docker rmi $(REPO)/$(DOCKER_IMAGE):$(VERSION)
	-rm -f $(DOCKER_IMAGE)

very-clean: clean-image
	-rm -f $(SOURCE_APP)
