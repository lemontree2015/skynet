#!/bin/sh

GOPATH=/root/mygo_projs # Note: Fix crontab env issue

# path
PROJ="monitor_server/gimcloud/skynet"
ROOT="`dirname "$0"`"
ROOT="`cd "$ROOT" && pwd`"

# default configurations
SKYNET_LOG_DIR="/data/skynet/log"
SKYNET_NOHUP_LOG=$SKYNET_LOG_DIR/skynet.NOHUP
SKYNET_DEFAULT_FLAGS="-stderrthreshold=INFO -log_dir=$SKYNET_LOG_DIR -conf_path=$ROOT/skynet.conf.default"

# skynet monitor server configurations
SKYNET_MONITOR_SERVER_CMD=$GOPATH/bin/skynet_monitor_server
SKYNET_MONITOR_SERVER_NOHUP_LOG=$SKYNET_LOG_DIR/skynet_monitor_server.NOHUP
SKYNET_MONITOR_SERVER_FLAGS="-stderrthreshold=INFO -log_dir=$SKYNET_LOG_DIR -conf_path=$ROOT/skynet.conf.default"

# benchmark/skynet_benchmark_service configurations
SKYNET_BENCHMARK_SERVICE_CMD=$GOPATH/bin/skynet_benchmark_service
SKYNET_BENCHMARK_SERVICE_NOHUP_LOG=$SKYNET_LOG_DIR/skynet_benchmark_service.NOHUP
SKYNET_BENCHMARK_SERVICE_FLAGS="-stderrthreshold=INFO -log_dir=$SKYNET_LOG_DIR -conf_path=$ROOT/skynet.conf.default"

# build *.go
run_build() {
  echo "start_run_build"

  go install $PROJ/client/conn
  echo "build $PROJ/client/conn success"
  go install $PROJ/client/loadbalancer
  echo "build $PROJ/client/loadbalancer success"
  go install $PROJ/client/pool
  echo "build $PROJ/client/pool success"
  go install $PROJ/client
  echo "build $PROJ/client success"
  go install $PROJ/config
  echo "build $PROJ/config success"
  go install $PROJ/cron
  echo "build $PROJ/cron success"
  go install $PROJ/global
  echo "build $PROJ/global success"
  go install $PROJ/misc
  echo "build $PROJ/misc success"
  go install $PROJ/rpc/bsonrpc
  echo "build $PROJ/rpc/bsonrpc success"
  go install $PROJ/service
  echo "build $PROJ/service success"
  go install $PROJ
  echo "build $PROJ success"

  # build skynet_monitor_server
  go install $PROJ/skynet_monitor_server/daemon
  echo "build $PROJ/skynet_monitor_server/daemon success"
  go install $PROJ/skynet_monitor_server
  echo "build $PROJ/skynet_monitor_server success"

  # build benchmark
  go install $PROJ/benchmark/skynet_benchmark_client
  echo "build $PROJ/benchmark/skynet_benchmark_client success"
  go install $PROJ/benchmark/skynet_benchmark_service
  echo "build $PROJ/benchmark/skynet_benchmark_service success"

  echo "end_run_build"
}


# run unit testcases
run_unit_test() {
  echo "start_run_unit_test"

  go test -v $PROJ $SKYNET_DEFAULT_FLAGS
  go test -v $PROJ/rpc/bsonrpc 
  go test -v $PROJ/client/pool $SKYNET_DEFAULT_FLAGS
  go test -v $PROJ/client/loadbalancer $SKYNET_DEFAULT_FLAGS
 
  echo "end_run_unit_test"
}

# run benchmark testcases
run_benchmark() {
  echo "start_run_benchmark"
  echo "end_run_benchmark"
}

usage() {
  echo "usage: skynet_ctl.sh"
  echo ""
  echo "The most commonly used skynet_ctl.sh commands are:"
  echo "  build                      build *.go"
  echo "  test                       run test/main.go"
  echo "  unit_test                  run unit testcases"
  echo "  benchmark                  run benchmark testcases"
  echo "  monitor_server_debug       start monitor_server in debug mode"
  echo "  monitor_server_nohup       start monitor_server in damon mode"
  echo "  benchmark_client_debug     start benchmark_client in debug mode"
  echo "  benchmark_service_debug    start benchmark_service in debug mode"
  echo "  benchmark_service_nohup    start benchmark_service in damon mode"
}


case "$1" in
  build)
    run_build
    ;;
  test)
     go run test/main.go $SKYNET_DEFAULT_FLAGS
    ;;
  unit_test)
    run_unit_test
    ;;
  benchmark)
    run_benchmark
    ;;
  monitor_server_debug)
    $SKYNET_MONITOR_SERVER_CMD $SKYNET_MONITOR_SERVER_FLAGS
    ;;
  monitor_server_nohup)
    nohup $SKYNET_MONITOR_SERVER_CMD $SKYNET_MONITOR_SERVER_FLAGS >$SKYNET_MONITOR_SERVER_NOHUP_LOG 2>&1 &
    ;;
  benchmark_client_debug)
    go run benchmark/skynet_benchmark_client/main.go $SKYNET_DEFAULT_FLAGS
    ;;
  benchmark_service_debug)
    $SKYNET_BENCHMARK_SERVICE_CMD $SKYNET_BENCHMARK_SERVICE_FLAGS
    ;;
  benchmark_service_nohup)
    nohup $SKYNET_BENCHMARK_SERVICE_CMD $SKYNET_BENCHMARK_SERVICE_FLAGS >$SKYNET_BENCHMARK_SERVICE_NOHUP_LOG 2>&1 &
    ;;
  *)
    usage
    ;;
esac
