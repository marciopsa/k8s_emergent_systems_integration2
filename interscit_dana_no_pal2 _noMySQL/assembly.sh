#!/bin/bash


#chown -R mysql:mysql /var/lib/mysql
#usermod -d /var/lib/mysql/ mysql
#/etc/init.d/mysql start


service mysql start

sleep 30


dana DCScheme.o

#dana -sp ../dc/ InteractiveAssembly.o ../../../dana/components/ws/core.o -p 2020
dana -sp ../dc/ InteractiveAssembly.o ../../dana/components/ws/core.o -p 2020
















#dana -sp ../dc/ EmergentSys.o ../../../dana/components/ws/core.o -p 2020


#docker build -t docker_interscity_env .
#docker run --name interscity_container -it docker_interscity_env bash

