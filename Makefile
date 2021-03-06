
.PHONY: repo-warning  
#
# This makefile assumes that docker is installed
#
# 3/14/2016 mln
#
# golog compile via docker container:
# https://hub.docker.com/_/golang/
#
### this needs to be evaluated. See make.golang
VERSION := 1.0.0

DOCKER_IMAGE := jsonsvalidator

#
# do some repo evals... want everyone to use their own, unless building for the team
#
DEFAULT_REPO := "quay.io/samsung_cnct"

# will need to be changed to samsung-cnct for prod version
DOCKER_REPO ?= DEFAULT_REPO
REPO := $(DOCKER_REPO)
#
#
#
# MAKE SUTE TO FILTER OUT THE CONTAINTER START...just in case
#
ARGLIST := $(filter-out golang-build-container,$(MAKECMDGOALS))
LOCALARGLIST := $(filter golang-build-container, $(MAKECMDGOALS))
MAKECMDGOALS := $(LOCALARGLIST)

repo-warning:
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
	@echo "vars:$(MAKE):$(MAKECMDGOALS):$(MAKEFLAGS):$(ARGLIST)"


#
# attempt to pass everything through to build script
# 
%::
	@echo "vars:$(MAKE):$(MAKECMDGOALS):$(MAKEFLAGS):$(ARGLIST)"
	./bin/build.sh -- $(MAKECMDGOALS) $(MAKEFLAGS) $(ARGLIST)

