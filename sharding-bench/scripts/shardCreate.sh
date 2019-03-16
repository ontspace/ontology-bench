#!/bin/bash


shardPerNode=2
servers=("10.1.4.4" "10.1.4.5")
#servers=("10.1.4.4" "10.1.4.5" "10.1.4.6" "10.1.4.7" "10.1.4.8" "10.1.4.9" "10.1.4.10" "10.1.4.11" "10.1.4.12" "10.1.4.13" "10.1.4.14" "10.1.4.15" "10.1.4.16" "10.1.4.17" "10.1.4.18" "10.1.4.19" "10.1.4.20" "10.1.4.21" "10.1.4.22" "10.1.4.23" "10.1.4.24" "10.1.4.26")


# create chains
id=1
counter=1
for s in "${servers[@]}"
do
	sed -i -e "s/wallet[0-9]*.dat/wallet$counter.dat/g" params/ShardCreate.json

	i=0
	while [ $i -lt $shardPerNode ]
	do
        echo "create shard " $id
		cat params/ShardCreate.json
		./bin/shardmgmt create config.json params/ShardCreate.json
        sleep 10
		((i++))
		((id++))
	done

	((counter++))
done

# config chains
id=1
counter=1
for s in "${servers[@]}"
do
	sed -i -e "s/wallet[0-9]*.dat/wallet$counter.dat/g" params/ShardConfig.json

	i=0
	while [ $i -lt $shardPerNode ]
	do
		sed -i -e "s/\"shard_id\":.*/\"shard_id\": $id,/g" params/ShardConfig.json
        echo "config shard " $id
		head -4 params/ShardConfig.json
		./bin/shardmgmt config config.json params/ShardConfig.json
        sleep 10
		((i++))
		((id++))
	done

	((counter++))
done

# peer join chain
id=1
counter=1
for s in "${servers[@]}"
do
	pk=`./bin/ontology account list -v -w wallets/wallet$counter.dat | grep 'Public key' | awk '{print $3}'`
	sed -i -e "s/wallet[0-9]*.dat/wallet$counter.dat/g" params/ShardPeerJoin.json
    sed -i -e "s/\"peer_pub_key\":.*/\"peer_pub_key\": \"$pk\",/g" params/ShardPeerJoin.json

	i=0
	while [ $i -lt $shardPerNode ]
	do
		sed -i -e "s/\"shard_id\":.*/\"shard_id\": $id,/g" params/ShardPeerJoin.json
        echo "peer join shard " $id
		head -5 params/ShardPeerJoin.json
		./bin/shardmgmt peerjoin config.json params/ShardPeerJoin.json
        sleep 10
		((i++))
		((id++))
	done

	((counter++))
done

# peer activate
id=1
counter=1
for s in "${servers[@]}"
do
	sed -i -e "s/wallet[0-9]*.dat/wallet$counter.dat/g" params/ShardActivate.json

	i=0
	while [ $i -lt $shardPerNode ]
	do
		sed -i -e "s/\"shard_id\":.*/\"shard_id\": $id/g" params/ShardActivate.json
        echo "peer activate " $id
		head -3 params/ShardActivate.json
		./bin/shardmgmt activate config.json params/ShardActivate.json
        sleep 10
		((i++))
		((id++))
	done

	((counter++))
done

