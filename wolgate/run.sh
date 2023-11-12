#!/bin/bash
CONFIG_PATH=/data/options.json
NEW_CONFIG="/config.json"

WOL_MAC=$(jq --raw-output '.wol_macAddress' $CONFIG_PATH)
WOL_BROADCAST=$(jq --raw-output '.wol_broadcastAddress' $CONFIG_PATH)
PROXY_HOST=$(jq --raw-output '.proxyServer_host' $CONFIG_PATH)
PROXY_PORT=$(jq --raw-output '.ports."25565/tcp"' $CONFIG_PATH)
DOMAINS=$(jq --raw-output '.domains' $CONFIG_PATH)

cat << EOF > $NEW_CONFIG
{
  "wol": {
    "macAddress": "$WOL_MAC",
    "broadcastAddress": "$WOL_BROADCAST"
  },
  "proxyServer": {
    "host": "$PROXY_HOST",
    "port": $PROXY_PORT
  },
  "domains": $DOMAINS
}
EOF

cat $NEW_CONFIG
/wolgate