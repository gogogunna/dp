#!/bin/bash

APPNAME="diplom"

kill() {
  docker ps --filter "ancestor=$APPNAME" -q | xargs docker kill
}

build() {
  docker build -t $APPNAME .
}

deploy() {
  kill
  build
  docker run -p 80:80 -d $APPNAME;
}

f_name=$1

"$f_name"