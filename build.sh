#!/bin/bash
#
# Script to start a docker golang container and run project Makefile.
#
# The current dictory is mounted to the docker container and 
# make -f make.golang is executed passing any additional arguments 
# that where supplied to this script.
#

SCRIPT_VERSION="4.0"

function gexample::build::debug {
    echo "[$(date)]: $*"
}
function gexample::build::info {
    echo "[$(date)]: $*"
}
function gexample::build::error {
    gexample::build::info $* >&2
}

function gexample::build::usage {
    gexample::build::error "Runs a golang build docker container and runs Makefile"
    gexample::build::error ""
    gexample::build::error "Usage:"
    gexample::build::error "   $0 [flags] -- [Makefile Args]"
    gexample::build::error ""
    gexample::build::error "Flags:"
    gexample::build::error "  -f,  --file :: golang make file name (make.golang)"
    gexample::build::error "  -h, -?, --help :: print usage"
    gexample::build::error "  -i, --int :: start an interactive shell"
    gexample::build::error "  -k, --kube :: route Makefile args to container build"
    gexample::build::error "  -m, --machine :: VM machine name, overrides DOCKER_MACHINE_NAME (gexample-build)"
    gexample::build::error "  -t, -test :: Test Docker Detection"
    gexample::build::error "  -v, --version :: print script verion"
    gexample::build::error "  -vv, --verbose :: more debug"
    gexample::build::error ""
    gexample::build::error "Env Vars:"
    gexample::build::error "  DOCKER_MACHINE_DRIVER :: (virtualbox) [optional]"
    gexample::build::error "  DOCKER_MACHINE_NAME :: (gexample-build) or set via argument [optional]"
    gexample::build::error ""
}
function gexample::build::version {
    gexample::build::info "$0 version ${SCRIPT_VERSION}"
}
#
# $1 = path to validate
function gexample::build::validate_tree {
  :
    #
    # validate the required source installation
    #
    #EXPECTED_BUILD_PATH="/src/github.com/samsung-cnct/golang-tools/example-project"
#    EXPECTED_BUILD_PATH="/src/github.com/samsung-cnct/jsonsvalidator"
#
#    if [ "${1}" != "${EXPECTED_BUILD_PATH}" ]; then
#        gexample::build::error "Expected build path ${EXPECTED_BUILD_PATH} not found."
#        gexample::build::error "Path ${1} found instead."
#        gexample::build::error "Your repo is not at the correct path."
#        gexample::build::error "See the README.md for the correct Directory Setup."
#        exit 2
#    else
#        gexample::build::info "Directory tree appears correct."
#    fi
}

# some best practice stuff
CRLF=$'\n'
CR=$'\r'
unset CDPATH

# XXX: this won't work if the last component is a symlink
[[ -z "$GOPATH" ]] && gexample::build::error "\$GOPATH is not set. Quitting"

[[ -d .git ]] && git_dir=$(pwd) || \
                 git_dir=$(find $(dirname $(pwd)) -type d -name .git)

go_dir="$GOPATH/src"
build_dir=$( echo ${git_dir#$go_dir})

#
gexample::build::info " "
gexample::build::info "=================================================="
gexample::build::info "   gexample golang build script via container"
gexample::build::info "   version: ${SCRIPT_VERSION}"
gexample::build::info "=================================================="
gexample::build::info "go_dir ${go_dir} "
gexample::build::info "git_dir ${git_dir} "
gexample::build::info "my_dir ${my_dir} "
gexample::build::info "build_dir ${build_dir} "
#
# validate the source tree setup
gexample::build::validate_tree ${build_dir}
#----------------------
# start the services first...this is so the ENV vars are available to the pods
#----------------------
#
# process args
#
MAKE_ARGS=""
VERBOSE=1
KUBE=0
INTERACTIVE=0
MAKEFILE_NAME="make.golang"
DOCKER_MACHINE_NAME=${DOCKER_MACHINE_NAME:-"gexample-build"}
TEST_DOCKER=0
# not set via args
readonly DOCKER_MACHINE_DRIVER=${DOCKER_MACHINE_DRIVER:-"virtualbox"}

while [ "$1" != "" ]; do
    case $1 in
        -v | --version )
            gexample::build::version
            exit
            ;;
        -h | -/? | --help )
            gexample::build::usage
            exit
            ;;
         -f | --file )
            shift
            MAKEFILE_NAME=$1
            ;; 
         -i | --int )
            INTERACTIVE=1
           ;;
         -m | --machine )
            shift
            DOCKER_MACHINE_NAME=$1
            ;; 
         -k | --kube )
            KUBE=1
           ;;
         -t | --test )
            TEST_DOCKER=1
           ;;
         -vv | --verbose )
            VERBOSE=1
           ;;
         -- ) 
            shift
            MAKE_ARGS=$@
            break
            ;;
         * )
             gexample::build::usage
             exit 1
    esac
    shift
done
MAKE_ARGS=$@

#gexample::build::info "VM: ${DOCKER_MACHINE_DRIVER} ${DOCKER_MACHINE_NAME} Verbose: $VERBOSE makefile: $MAKEFILE_NAME make_args: ${MAKE_ARGS} kubemake: $KUBE" 

#function gexample::build::machinecheck {
#    #
#    # check if docker is running locally already first (e.g. docker for os-x)
#    # if not, then check for docker-machine
#    # if not, then ask to install docker locallly or docker-machine
#    #
#    if [[ -z "$(which docker)" ]]; then
#        if [[ -z "$(which docker-machine)" ]]; then
#            gexample::build::info "Neither docker nor docker-machine is not found... please install one of them."
#            exit 1
#        elif [[ -n "$(which docker-machine)" ]]; then
#            gexample::build::info "docker-machine was found"
#            docker-machine inspect  "${DOCKER_MACHINE_NAME}" > /dev/null || {
#                gexample::build::info "Creating a docker-machine instance for build: ${DOCKER_MACHINE_NAME}"
#                docker-machine create --driver "${DOCKER_MACHINE_DRIVER}" "${DOCKER_MACHINE_NAME}" > /dev/null || {
#                gexample::build::error "Something went wrong creating a machine."
#                gexample::build::error "Try the following: "
#                gexample::build::error "docker-machine create -d ${DOCKER_MACHINE_DRIVER} ${DOCKER_MACHINE_NAME}"
#                return 1
#                } 
#            }
#        fi
#    else
#        DUMMY=$(docker info 2>/dev/null)
#        if [ $? -ne 0 ]; then
#            gexample::build::info "Docker is installed by not running.  Please run docker. And try your command again"
#            exit 1
#        else
#            # docker is running...we can use it
#            gexample::build::info "Docker is running, continuing."
#        fi
#    fi
#}
#
#gexample::build::machinecheck
#if [ $TEST_DOCKER == 1 ];then
#    gexample::build::info "Only Tested for Docker due to -t flag."
#    exit 0
#fi
#
# run the Makefile to build
#
#
#  Use our custom container that has golang and glide installed.
#
#GOLANG_CONTAINER="golang:1.6"
#  Make settable via args/env vars
GOLANG_CONTAINER=${GOLANG_CONTAINER:-"quay.io/samsung_cnct/goglide:0.1.8"}

BUILD_VERSION="0.0.3"

function gexample::build::interactive {
    gexample::build::info "Running Interactive ${GOLANG_CONTAINER}"
    docker run \
        --rm \
        -it \
        --name golang-build-container \
        -v ${go_dir}:/go \
        -w /go${build_dir} \
        -e VERSION=${BUILD_VERSION} \
        -e LOCAL_USER=$USER \
        ${GOLANG_CONTAINER} \
        bash
}

#function gexample::build::make::orig {
#    gexample::build::info "Running Makefile: ${MAKEFILE_NAME} in ${GOLANG_CONTAINER}"
#    docker run \
#        --rm \
#        --name golang-build-container \
#        -v ${go_dir}:/go \
#        -w /go${build_dir} \
#        -e VERSION=${BUILD_VERSION} \
#        -e LOCAL_USER=$USER \
#        ${GOLANG_CONTAINER} \
#        bash -c "pwd;\
#        df;\
#        ls -l;\
#        env|sort;\
#        which make;\
#        whoami; \
#        make --version;\
#        make --file ${MAKEFILE_NAME} ${MAKE_ARGS};"
#}

function gexample::build::make()
{
    make -C make.golang ${MAKE_ARGS}
}

function gexample::build::container {
    gexample::build::info "Running Makefile: ${MAKEFILE_NAME} in current shell"
    make -C ./_containerize  ${MAKE_ARGS}
}


if [ $INTERACTIVE == 1 ]; then
    gexample::build::info "-------------------------------------"
    gexample::build::info "        Interactive golang"
    gexample::build::info "-------------------------------------"
    gexample::build::interactive
else
if [ $KUBE == 0 ]; then
    gexample::build::info "-------------------------------------"
    gexample::build::info "        Building golang app"
    gexample::build::info "-------------------------------------"
    gexample::build::make
else
    gexample::build::info "*************************************"
    gexample::build::info "Building docker object for kubernetes"
    gexample::build::info "*************************************"
    gexample::build::container
fi
fi


gexample::build::info "------------ build script finished --------------"
exit
