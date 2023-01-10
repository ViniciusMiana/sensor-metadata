#!/bin/bash

set -e


usage() {
  echo "Usage: $0 [START|STOP]"
}

start_mongo() {
  echo "[STATUS] Starting Mongo DB"
  docker run --rm -d -p 27017:27017 --name mongodb-sensor-test mongo:4.2.23 --replSet sensor0 > /dev/null
  while ! curl -s localhost:27017/ > /dev/null;
    do sleep 1
  done
  docker exec -it mongodb-sensor-test mongo --eval 'rs.initiate({"_id":"sensor0","members":[{"_id":0,"host":"localhost:27017"}]})'
  echo "[STATUS] Mongo DB Started"
}

stop_mongo() {
  if [ "`docker ps -q -f name=mongodb-sensor-test`" != "" ]; then
    docker container stop mongodb-sensor-test > /dev/null
    echo "[STATUS] Mongo DB stopped"
  fi
}

if [ "$1" != "START" ] && [ "$1" != "STOP" ]; then
  usage
  exit 1
fi

if [ "`command -v docker`" == "" ]; then
  echo "[ERROR] Cannot start the container. Docker is not installed on your machine. Please install docker and try again"
  exit 1
fi


stop_mongo
if [ "$1" == "START" ]; then
    start_mongo
fi
