#!/bin/bash
while true; do
  zt agent list | grep -o -e 'host[a-zA-z0-9\-]*' | xargs -I{} sh -c 'zt agent stack --app-alias $1 > $1.$(date +%s).stack' -- {}
  sleep 15
done
