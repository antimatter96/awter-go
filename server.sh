#!/usr/bin/env bash
tput setab 0;

tput setaf 3;
echo "Building"
tput setaf 7;
go build -o awter-go -race -ldflags "-X 'main.build=$(date +%Y-%m-%d_%H:%M)'"
tput setaf 2;
echo "Built"
tput setaf 7;
# ./awter-go --port=9005 --mysqlURL="root:661996@/awter_db_test?parseTime=true"
./awter-go --port 9005 --mysqlURL ""
