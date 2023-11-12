#!/bin/bash
CONFIG_PATH=/data/options.json

# Exemple de lecture des options
WOL_MAC_ADDRESS=$(jq --raw-output '.wol_mac_address' $CONFIG_PATH)
WOL_BROADCAST_ADDRESS=$(jq --raw-output '.wol_broadcast_address' $CONFIG_PATH)
PROXY_SERVER_HOST=$(jq --raw-output '.proxy_server_host' $CONFIG_PATH)
PROXY_SERVER_PORT=$(jq --raw-output '.proxy_server_port' $CONFIG_PATH)

# Vous pouvez ensuite utiliser ces variables pour configurer votre application
# Cela peut impliquer de générer un fichier de configuration ou de passer des arguments à votre application

# Lancer votre application
exec /path/to/your/app --wol_mac_address "$WOL_MAC_ADDRESS" --wol_broadcast_address "$WOL_BROADCAST_ADDRESS" --proxy_server_host "$PROXY_SERVER_HOST" --proxy_server_port "$PROXY_SERVER_PORT"
