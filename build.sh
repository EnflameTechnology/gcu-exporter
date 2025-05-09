#!/bin/bash
#
# Copyright 2023 Enflame. All Rights Reserved.
#
set -eu -o pipefail

function usage() {
  cat <<EOF
Usage: ./build.sh [command]
Command:
    all     Build all bins and make tar package.
    clean   Clean all bins and make tar package.
Example:
    ./build.sh all
    ./build.sh clean
EOF
}

PKG_VER="1.0.0"

function clean_all() {
    make clean
    if [ -d "dist" ]; then
        rm -rf dist
    fi
}
function build_all() {
    echo -e "\033[33mBuilding gcu-exporter package.\033[0m"

    make #>/dev/null 2>&1

    if [ -d "dist" ]; then
        rm -rf dist
    fi

    PKG_DIR="dist/gcu-exporter_${PKG_VER}"
    mkdir -p ${PKG_DIR}

    cp -rf -L deployments/* ${PKG_DIR}
    cp -rf LICENSE ${PKG_DIR}
    mv gcu-exporter ${PKG_DIR}

    echo -e "\033[33mgcu-exporter released successfully.\033[0m"
}

function main() {
    [[ "$EUID" -ne 0 ]] && { echo "[INFO] Run as non-root"; }
    [[ "$#" -lt 1 ]] && { usage >&2; exit 1; }

    case $1 in
    "clean")
        clean_all
        ;;
    "all")
        build_all
       ;;
    *)
        usage
        exit 1
        ;;
    esac
}

main "$@"
