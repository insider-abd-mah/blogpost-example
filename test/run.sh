#!/bin/bash

SQL_PATH="$1"
COMPLETED_INFORMATION="completed"

function sendSQLFileToEntryPoint() {
  docker cp "$SQL_PATH" mysql_fake:/docker-entrypoint-initdb.d/init.sql
  echo $COMPLETED_INFORMATION
}

function newContainer() {
  docker kill mysql_fake
  docker container rm -f mysql_fake
  docker run -d --name mysql_fake -p 3305:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=sample_db mariadb:latest
  sendSQLFileToEntryPoint
}

TEST_CONTAINER_STATUS="$(docker inspect -f '{{.State.Running}}' mysql_fake)"

if [ "$TEST_CONTAINER_STATUS" != "true" ]; then
    newContainer
fi

echo $COMPLETED_INFORMATION
