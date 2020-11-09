#!/bin/bash

set -euo pipefail

NETWORK=${1:-testnet}

docker build . -t proxy
docker run -it --rm --name proxy \
-e "NETWORK=$NETWORK" \
--net "${NETWORK}_default" \
-p 8080:8080 \
-v /var/run/docker.sock:/var/run/docker.sock \
-v "$HOME/.xud-docker/$NETWORK/data/xud:/root/.xud" \
-v "$HOME/.xud-docker/$NETWORK/data/lndbtc:/root/.lndbtc" \
-v "$HOME/.xud-docker/$NETWORK/data/lndltc:/root/.lndltc" \
-v "$HOME/.xud-docker/$NETWORK/logs/config.sh:/root/config.sh" \
-v "$HOME/.xud-docker/$NETWORK:/root/network:ro" \
proxy
