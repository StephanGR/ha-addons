#!/bin/bash
CONFIG_PATH="/data/options.json"
NEW_CONFIG="/config.json"

DOMAINS=$(jq --raw-output '.domains' $CONFIG_PATH)

cat << EOF > $NEW_CONFIG
{
  "domains": $DOMAINS
}
EOF

echo "Loaded configuration :"
jq '.' $NEW_CONFIG

/wolgate