// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by 'make imports' - DO NOT EDIT.

/*
Package include imports all Module and MetricSet packages so that they register
their factories with the global registry. This package can be imported in the
main package to automatically register all of the standard supported Metricbeat
modules.
*/
package include

import (
	_ "github.com/elastic/beats/metricbeat/module/aerospike"
	_ "github.com/elastic/beats/metricbeat/module/aerospike/namespace"
	_ "github.com/elastic/beats/metricbeat/module/apache"
	_ "github.com/elastic/beats/metricbeat/module/apache/status"
	_ "github.com/elastic/beats/metricbeat/module/beat"
	_ "github.com/elastic/beats/metricbeat/module/beat/state"
	_ "github.com/elastic/beats/metricbeat/module/beat/stats"
	_ "github.com/elastic/beats/metricbeat/module/ceph"
	_ "github.com/elastic/beats/metricbeat/module/ceph/cluster_disk"
	_ "github.com/elastic/beats/metricbeat/module/ceph/cluster_health"
	_ "github.com/elastic/beats/metricbeat/module/ceph/cluster_status"
	_ "github.com/elastic/beats/metricbeat/module/ceph/monitor_health"
	_ "github.com/elastic/beats/metricbeat/module/ceph/osd_df"
	_ "github.com/elastic/beats/metricbeat/module/ceph/osd_tree"
	_ "github.com/elastic/beats/metricbeat/module/ceph/pool_disk"
	_ "github.com/elastic/beats/metricbeat/module/consul"
	_ "github.com/elastic/beats/metricbeat/module/consul/agent"
	_ "github.com/elastic/beats/metricbeat/module/couchbase"
	_ "github.com/elastic/beats/metricbeat/module/couchbase/bucket"
	_ "github.com/elastic/beats/metricbeat/module/couchbase/cluster"
	_ "github.com/elastic/beats/metricbeat/module/couchbase/node"
	_ "github.com/elastic/beats/metricbeat/module/couchdb"
	_ "github.com/elastic/beats/metricbeat/module/couchdb/server"
	_ "github.com/elastic/beats/metricbeat/module/docker"
	_ "github.com/elastic/beats/metricbeat/module/docker/container"
	_ "github.com/elastic/beats/metricbeat/module/docker/cpu"
	_ "github.com/elastic/beats/metricbeat/module/docker/diskio"
	_ "github.com/elastic/beats/metricbeat/module/docker/event"
	_ "github.com/elastic/beats/metricbeat/module/docker/healthcheck"
	_ "github.com/elastic/beats/metricbeat/module/docker/image"
	_ "github.com/elastic/beats/metricbeat/module/docker/info"
	_ "github.com/elastic/beats/metricbeat/module/docker/memory"
	_ "github.com/elastic/beats/metricbeat/module/docker/network"
	_ "github.com/elastic/beats/metricbeat/module/dropwizard"
	_ "github.com/elastic/beats/metricbeat/module/dropwizard/collector"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/ccr"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/cluster_stats"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index_recovery"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index_summary"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/ml_job"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/node"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/node_stats"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/pending_tasks"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/shard"
	_ "github.com/elastic/beats/metricbeat/module/envoyproxy"
	_ "github.com/elastic/beats/metricbeat/module/envoyproxy/server"
	_ "github.com/elastic/beats/metricbeat/module/etcd"
	_ "github.com/elastic/beats/metricbeat/module/etcd/leader"
	_ "github.com/elastic/beats/metricbeat/module/etcd/metrics"
	_ "github.com/elastic/beats/metricbeat/module/etcd/self"
	_ "github.com/elastic/beats/metricbeat/module/etcd/store"
	_ "github.com/elastic/beats/metricbeat/module/golang"
	_ "github.com/elastic/beats/metricbeat/module/golang/expvar"
	_ "github.com/elastic/beats/metricbeat/module/golang/heap"
	_ "github.com/elastic/beats/metricbeat/module/graphite"
	_ "github.com/elastic/beats/metricbeat/module/graphite/server"
	_ "github.com/elastic/beats/metricbeat/module/haproxy"
	_ "github.com/elastic/beats/metricbeat/module/haproxy/info"
	_ "github.com/elastic/beats/metricbeat/module/haproxy/stat"
	_ "github.com/elastic/beats/metricbeat/module/http"
	_ "github.com/elastic/beats/metricbeat/module/http/json"
	_ "github.com/elastic/beats/metricbeat/module/http/server"
	_ "github.com/elastic/beats/metricbeat/module/jolokia"
	_ "github.com/elastic/beats/metricbeat/module/jolokia/jmx"
	_ "github.com/elastic/beats/metricbeat/module/kafka"
	_ "github.com/elastic/beats/metricbeat/module/kafka/consumergroup"
	_ "github.com/elastic/beats/metricbeat/module/kafka/partition"
	_ "github.com/elastic/beats/metricbeat/module/kibana"
	_ "github.com/elastic/beats/metricbeat/module/kibana/stats"
	_ "github.com/elastic/beats/metricbeat/module/kibana/status"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/apiserver"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/container"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/controllermanager"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/event"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/node"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/pod"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/proxy"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_container"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_deployment"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_node"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_pod"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_replicaset"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/state_statefulset"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/system"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/util"
	_ "github.com/elastic/beats/metricbeat/module/kubernetes/volume"
	_ "github.com/elastic/beats/metricbeat/module/kvm"
	_ "github.com/elastic/beats/metricbeat/module/kvm/dommemstat"
	_ "github.com/elastic/beats/metricbeat/module/logstash"
	_ "github.com/elastic/beats/metricbeat/module/logstash/node"
	_ "github.com/elastic/beats/metricbeat/module/logstash/node_stats"
	_ "github.com/elastic/beats/metricbeat/module/memcached"
	_ "github.com/elastic/beats/metricbeat/module/memcached/stats"
	_ "github.com/elastic/beats/metricbeat/module/mongodb"
	_ "github.com/elastic/beats/metricbeat/module/mongodb/collstats"
	_ "github.com/elastic/beats/metricbeat/module/mongodb/dbstats"
	_ "github.com/elastic/beats/metricbeat/module/mongodb/metrics"
	_ "github.com/elastic/beats/metricbeat/module/mongodb/replstatus"
	_ "github.com/elastic/beats/metricbeat/module/mongodb/status"
	_ "github.com/elastic/beats/metricbeat/module/munin"
	_ "github.com/elastic/beats/metricbeat/module/munin/node"
	_ "github.com/elastic/beats/metricbeat/module/mysql"
	_ "github.com/elastic/beats/metricbeat/module/mysql/galera_status"
	_ "github.com/elastic/beats/metricbeat/module/mysql/status"
	_ "github.com/elastic/beats/metricbeat/module/nats"
	_ "github.com/elastic/beats/metricbeat/module/nats/connections"
	_ "github.com/elastic/beats/metricbeat/module/nats/routes"
	_ "github.com/elastic/beats/metricbeat/module/nats/stats"
	_ "github.com/elastic/beats/metricbeat/module/nats/subscriptions"
	_ "github.com/elastic/beats/metricbeat/module/nginx"
	_ "github.com/elastic/beats/metricbeat/module/nginx/stubstatus"
	_ "github.com/elastic/beats/metricbeat/module/php_fpm"
	_ "github.com/elastic/beats/metricbeat/module/php_fpm/pool"
	_ "github.com/elastic/beats/metricbeat/module/php_fpm/process"
	_ "github.com/elastic/beats/metricbeat/module/postgresql"
	_ "github.com/elastic/beats/metricbeat/module/postgresql/activity"
	_ "github.com/elastic/beats/metricbeat/module/postgresql/bgwriter"
	_ "github.com/elastic/beats/metricbeat/module/postgresql/database"
	_ "github.com/elastic/beats/metricbeat/module/postgresql/statement"
	_ "github.com/elastic/beats/metricbeat/module/prometheus"
	_ "github.com/elastic/beats/metricbeat/module/prometheus/collector"
	_ "github.com/elastic/beats/metricbeat/module/rabbitmq"
	_ "github.com/elastic/beats/metricbeat/module/rabbitmq/connection"
	_ "github.com/elastic/beats/metricbeat/module/rabbitmq/exchange"
	_ "github.com/elastic/beats/metricbeat/module/rabbitmq/node"
	_ "github.com/elastic/beats/metricbeat/module/rabbitmq/queue"
	_ "github.com/elastic/beats/metricbeat/module/redis"
	_ "github.com/elastic/beats/metricbeat/module/redis/info"
	_ "github.com/elastic/beats/metricbeat/module/redis/key"
	_ "github.com/elastic/beats/metricbeat/module/redis/keyspace"
	_ "github.com/elastic/beats/metricbeat/module/system"
	_ "github.com/elastic/beats/metricbeat/module/system/core"
	_ "github.com/elastic/beats/metricbeat/module/system/cpu"
	_ "github.com/elastic/beats/metricbeat/module/system/diskio"
	_ "github.com/elastic/beats/metricbeat/module/system/entropy"
	_ "github.com/elastic/beats/metricbeat/module/system/filesystem"
	_ "github.com/elastic/beats/metricbeat/module/system/fsstat"
	_ "github.com/elastic/beats/metricbeat/module/system/load"
	_ "github.com/elastic/beats/metricbeat/module/system/memory"
	_ "github.com/elastic/beats/metricbeat/module/system/network"
	_ "github.com/elastic/beats/metricbeat/module/system/process"
	_ "github.com/elastic/beats/metricbeat/module/system/process_summary"
	_ "github.com/elastic/beats/metricbeat/module/system/raid"
	_ "github.com/elastic/beats/metricbeat/module/system/socket"
	_ "github.com/elastic/beats/metricbeat/module/system/socket_summary"
	_ "github.com/elastic/beats/metricbeat/module/system/uptime"
	_ "github.com/elastic/beats/metricbeat/module/traefik"
	_ "github.com/elastic/beats/metricbeat/module/traefik/health"
	_ "github.com/elastic/beats/metricbeat/module/uwsgi"
	_ "github.com/elastic/beats/metricbeat/module/uwsgi/status"
	_ "github.com/elastic/beats/metricbeat/module/vsphere"
	_ "github.com/elastic/beats/metricbeat/module/vsphere/datastore"
	_ "github.com/elastic/beats/metricbeat/module/vsphere/host"
	_ "github.com/elastic/beats/metricbeat/module/vsphere/virtualmachine"
	_ "github.com/elastic/beats/metricbeat/module/windows"
	_ "github.com/elastic/beats/metricbeat/module/windows/perfmon"
	_ "github.com/elastic/beats/metricbeat/module/windows/service"
	_ "github.com/elastic/beats/metricbeat/module/zookeeper"
	_ "github.com/elastic/beats/metricbeat/module/zookeeper/connection"
	_ "github.com/elastic/beats/metricbeat/module/zookeeper/mntr"
	_ "github.com/elastic/beats/metricbeat/module/zookeeper/server"
)
