FROM centos:7

ARG BOSH_VERSION=6.2.1
ARG BOSH_DEPLOYMENT_VERSION=c546c5110124954271e9c9d6c41775a0a18bd8d7
ARG GARDEN_VERSION=1.19.12

RUN yum install -y \
    # Tools
    wget unzip git sudo xfsprogs \
    # BOSH dependencies (from https://bosh.io/docs/cli-v2-install/)
    gcc gcc-c++ ruby ruby-devel mysql-devel postgresql-devel postgresql-libs sqlite-devel libxslt-devel libxml2-devel patch openssl gem install yajl-ruby

RUN wget -q -O /usr/local/bin/bosh https://github.com/cloudfoundry/bosh-cli/releases/download/v${BOSH_VERSION}/bosh-cli-${BOSH_VERSION}-linux-amd64 && chmod +x /usr/local/bin/bosh
RUN wget -q -O /usr/local/bin/gdn wget https://github.com/cloudfoundry/garden-runc-release/releases/download/v${GARDEN_VERSION}/gdn-${GARDEN_VERSION}

WORKDIR /root
RUN git clone https://github.com/cloudfoundry/bosh-deployment.git \
  && cd bosh-deployment \
  && git checkout -q -f ${BOSH_DEPLOYMENT_VERSION}

COPY create_env.sh /usr/local/bin/create_env.sh

# This fails on docker build, but pre-downloads dependencies
#RUN create_env.sh || true
