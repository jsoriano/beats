#!/bin/bash

CREDS_PATH=/root/creds.yml
BOSH_DEPLOYMENT=/root/bosh-deployment

bosh create-env ${BOSH_DEPLOYMENT}/bosh.yml \
  --ops-file ${BOSH_DEPLOYMENT}/docker/cpi.yml \
  --ops-file ${BOSH_DEPLOYMENT}/docker/unix-sock.yml \
  --ops-file ${BOSH_DEPLOYMENT}/uaa.yml \
  --ops-file ${BOSH_DEPLOYMENT}/credhub.yml \
  --ops-file ${BOSH_DEPLOYMENT}/jumpbox-user.yml \
  --vars-store ${CREDS_PATH} \
  --var director_name=bosh-lite \
  --var docker_host=unix:///var/run/docker.sock \
  --var network=bosh \
  --var internal_cidr=10.0.0.0/24 \
  --var internal_gw=10.0.0.1 \
  --var internal_ip=10.0.0.100

export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh int ./creds.yml --path /admin_password)
bosh alias-env local -e 10.0.0.100 --ca-cert <(bosh int ./creds.yml --path /director_ssl/ca)


