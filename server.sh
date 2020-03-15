#!/usr/bin/env bash

go build
./awter-go --port=9005 --mysqlURL="root:661996@/awter_db_test?parseTime=true"

