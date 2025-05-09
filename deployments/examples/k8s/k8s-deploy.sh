#!/bin/bash
#--------------------------------------------------
# This script is gcu-exporter examples.
# @usage: ./k8s-deploy.sh --help
#--------------------------------------------------
set -eu -o pipefail

function usage() {
  cat <<EOF
Usage: k8s-deploy.sh [command]

Commands:
    apply    Apply k8s yaml
    delete   Delete k8s yaml

Examples:
    k8s-deploy.sh apply
    k8s-deploy.sh delete
EOF
}

function apply() {
    echo -e "[INFO] \033[33mApply k8s yaml\033[0m"
    if [ ! -d /data ]; then
        mkdir -p /data
    fi
    if [ ! -d /data/prometheus ]; then
        mkdir -p /data/prometheus
    fi
    kubectl apply -f yaml/namespace.yaml
    kubectl apply -f yaml/gcu-exporter.yaml
    kubectl apply -f yaml/prometheus.yaml
    kubectl apply -f yaml/grafana.yaml
}

function delete() {
    echo -e "[INFO] \033[33mDelete k8s yaml\033[0m"
    kubectl delete -f yaml/prometheus.yaml
    kubectl delete -f yaml/grafana.yaml
    kubectl delete -f yaml/gcu-exporter.yaml
    kubectl delete -f yaml/namespace.yaml
}

function main() {

    [[ "$EUID" -ne 0 ]] && { echo "[ERROR] Must run as root"; exit 1; }
    [[ "$#" -eq 0 ]] && { usage >&2; exit 1; }

    ACTION="*"
    case "$1" in
        apply)
            ACTION="apply"
            ;;
        delete)
            ACTION="delete"
            ;;
        *)
            usage
            exit 1
            ;;
    esac

    [[ "$ACTION" == "" ]] && { echo "[ERROR] illegal option"; usage; exit 1; }

    echo -e "[INFO] \033[33mAction start\033[0m : $ACTION"
    ${ACTION} || { echo -e "[ERROR] \033[31mAction is failed\033[0m : $ACTION"; return 1; }
    echo -e "[INFO] \033[32mAction is successful.\033[0m : $ACTION"
}

main "$@"
