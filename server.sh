#!/usr/bin/env bash
echo "Building"
go build -o awter-go -race -ldflags "-X 'main.build=$(date +%Y-%m-%d\ %H:%M)'"
echo "Built"
# ./awter-go --port=9005 --mysqlURL="root:661996@/awter_db_test?parseTime=true"
./awter-go --port=9005 --mysqlURL=""
