#!/bin/bash

cd "$(dirname "$0")"

git checkout main && \
git pull origin main && \
make -f Docker.mk build && \
make -f Docker.mk stop
make -f Docker.mk start