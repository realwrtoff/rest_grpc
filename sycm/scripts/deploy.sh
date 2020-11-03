#!/usr/bin/env bash

source tmp/environment.sh
Version=$(git describe --tags)

function Trac() {
    echo "[TRAC] [$(date +"%Y-%m-%d %H:%M:%S")] $1"
}

function Info() {
    echo "\033[1;32m[INFO] [$(date +"%Y-%m-%d %H:%M:%S")] $1\033[0m"
}

function Warn() {
    echo "\033[1;31m[WARN] [$(date +"%Y-%m-%d %H:%M:%S")] $1\033[0m"
    return 1
}

function Build() {
    cd .. && make image && cd -
    docker login --username="${RegistryUsername}" --password="${RegistryPassword}" "${RegistryServer}"
    docker tag "${Image}:${Version}" "${RegistryServer}/${Image}:${Version}"
    docker push "${RegistryServer}/${Image}:${Version}"
}

function Help() {
    echo "sh deploy.sh <action>"
    echo "example"
    echo "  sh deploy.sh build"
}

function main() {
    if [ -z "$1" ]; then
        Help
        echo $Version
        return 0
    fi

    case "$1" in
        "build") Build;;
    esac
}

main "$@"
