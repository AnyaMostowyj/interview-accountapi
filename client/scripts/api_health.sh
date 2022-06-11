#!/usr/bin/env bash

while [ "$(curl --connect-timeout 2 -s -o /dev/null -w '%{http_code}' accountapi:8080/v1/health)" != "200" ]
do 
    echo checking api health
    sleep 5
done

echo api is running
go test -v -host "http://accountapi:8080"
