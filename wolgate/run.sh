#!/bin/bash
CONFIG_PATH="/data/options.json"
NEW_CONFIG="/config.json"

WOL_MAC=$(jq --raw-output '.wol_macAddress' $CONFIG_PATH)
WOL_BROADCAST=$(jq --raw-output '.wol_broadcastAddress' $CONFIG_PATH)
DOMAINS=$(jq --raw-output '.domains' $CONFIG_PATH)

cat << EOF > $NEW_CONFIG
{
  "wol": {
    "macAddress": "$WOL_MAC",
    "broadcastAddress": "$WOL_BROADCAST"
  },
  "domains": $DOMAINS
}
EOF

echo "Loaded configuration :"
jq '.' $NEW_CONFIG

/wolgate