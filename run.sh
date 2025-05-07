#!/bin/bash

APPNAME="diplom"

kill() {
  docker ps --filter "ancestor=$APPNAME" -q | xargs docker kill
}

build() {
  docker build --no-cache -t $APPNAME --file Dockerfile --build-context certs_context=/Users/ivanbrynov/GolandProjects/diplomca .
}

build_prod() {
  docker build -t $APPNAME --file Dockerfile --build-context certs_context=/etc/letsencrypt/archive/tinvestanalytics.ru .
}

deploy() {
  kill
  build
  docker run -p 443:443 -d $APPNAME;
}

deploy_prod() {
  kill
  build_prod
  docker run -p 443:443 -d $APPNAME;
}

f_name=$1

"$f_name"