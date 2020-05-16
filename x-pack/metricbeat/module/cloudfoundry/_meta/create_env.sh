#!/bin/bash

CREDS_PATH=/root/creds.yml
BOSH_DEPLOYMENT=/root/bosh-deployment

bosh create-env ${BOSH_DEPLOYMENT}/bosh.yml \
  --ops-file ${BOSH_DEPLOYMENT}/bosh-lite.yml \
  --ops-file ${BOSH_DEPLOYMENT}/warden/cpi.yml \
  --ops-file ${BOSH_DEPLOYMENT}/uaa.yml \
  --ops-file ${BOSH_DEPLOYMENT}/credhub.yml \
  --ops-file ${BOSH_DEPLOYMENT}/jumpbox-user.yml \
  --vars-store ${CREDS_PATH} \
  --var director_name=bosh-lite \
  --var internal_cidr=10.0.0.0/24 \
  --var internal_gw=10.0.0.1 \
  --var internal_ip=10.0.0.100 \
  --var garden_host=/var/run/docker.sock
