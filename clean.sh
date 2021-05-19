#!/bin/bash

containers=$(docker ps -a -q)
volumes=$(docker volume ls -q)

[ -z "$containers" ] && echo "No Containers to Remove" || docker rm -f $containers
[ -z "$volumes" ] && echo "No Volumes to Remove" || docker volume rm $volumes
fuser -k 8020/tcp || echo "Nothing on Port 8020"
fuser -k 8040/tcp || echo "Nothing on Port 8040"
fuser -k 8060/tcp || echo "Nothing on Port 8060"
fuser -k 8080/tcp || echo "Nothing on Port 8080"
fuser -k 8090/tcp|| echo "Nothing on Port 8090"