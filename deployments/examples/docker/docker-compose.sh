#!/bin/bash
#-----------------------------------------------------------
# This script is gcu-exporter examples
#------------------------------------------------------------
set -eu -o pipefail

function usage() {
  cat <<EOF
Usage: $0 [command]

Commands:
    init        init docker compose
    up          docker compose up -d
    down        docker compose down
Examples:
    $0 up
    $0 down
EOF
} 

function init() {
    echo -e "[INFO] \033[33mdocker-compose init\033[0m"

    DIST=$(cat /etc/issue | awk '{print $1}')
    if [ "Ubuntu" == "${DIST}" ]; then
        which docker-compose &> /dev/null
        if [ $? != 0 ]; then
            apt install -y docker-compose
        fi
    elif [ -f /etc/centos-release ]; then
        which docker-compose &> /dev/null
        if [ $? != 0 ]; then
            yum install -y docker-compose
        fi
    else
        echo "[ERROR] ${OSDist} not supported!"; exit 1;
    fi

    echo -e "[INFO] \033[33mpull docker images....\033[0m"

    docker pull prom/prometheus:latest
    # docker pull prom/alertmanager:latest
    docker pull grafana/grafana:latest

    # pull from internel docker hub, or build the image manually
    docker pull artifact.enflame.cn/enflame_docker_images/enflame/gcu-exporter:latest

    if [ ! -d /var/lib/grafana ]; then
        mkdir -p /var/lib/grafana
    fi
    chmod 777 -R /var/lib/grafana
}

function up() {
    echo -e "[INFO] \033[33mdocker-compose up -d\033[0m"
    if [ ! -d /var/lib/grafana ]; then
        mkdir -p /var/lib/grafana
    fi
    chmod 777 -R /var/lib/grafana

    docker-compose -f ./docker-compose.yaml up -d
}

function down() {
    echo -e "[INFO] \033[33mdocker-compose down\033[0m"
    docker-compose -f ./docker-compose.yaml down
}

function main() {
    [[ "$EUID" -ne 0 ]] && { echo "[ERROR] Must run as root"; exit 1; }
    [[ "$#" -eq 0 ]] && { usage >&2; exit 1; }

    ACTION="*"
    case "$1" in
        init)
            ACTION="init"
            ;;
        up)
            ACTION="up"
            ;;
        down)
            ACTION="down"
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
