#!/bin/bash

rootServer="10.1.4.27"
echo $rootServer
ssh  $rootServer 'cd /mnt/shard && ./ontology info curblockheight'

shardPerNode=2
servers=("10.1.4.4" "10.1.4.5")
#servers=("10.1.4.4" "10.1.4.5" "10.1.4.6" "10.1.4.7" "10.1.4.8" "10.1.4.9" "10.1.4.10" "10.1.4.11" "10.1.4.12" "10.1.4.13" "10.1.4.14" "10.1.4.15" "10.1.4.16" "10.1.4.17" "10.1.4.18" "10.1.4.19" "10.1.4.20" "10.1.4.21" "10.1.4.22" "10.1.4.23" "10.1.4.24" "10.1.4.26")

id=1
for s in "${servers[@]}"
do
    i=0
    while [ $i -lt $shardPerNode ]
    do
        echo $s '--' $id
        port=$((20336 + id * 10))
        ssh  $s "cd /mnt/shard && ./ontology info curblockheight --rpcport $port"
        ((id++))
        ((i++))
    done
done

