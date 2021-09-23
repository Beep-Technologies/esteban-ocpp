#!/usr/bin/env bash
if ! command -v swag &> /dev/null
then
    echo -e "swag cli not found. install swag \n\tgithub.com/swaggo/swag/cmd/swag\n and run this script again"
    exit 1
fi

swag init -d cmd/bb3-ocpp-ws,api/rpc,api/rest/controller -g main.go
