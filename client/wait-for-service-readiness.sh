#!/usr/bin/env bash

while [[ "$$(curl --connect-timeout 2 -s -o /dev/null -w ''%{http_code}'' accountapi:8080/v1/health)" != "200" ]]; do echo ..; sleep 5; done; echo backend is up
