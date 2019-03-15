#!/bin/bash

walletCnt=30
counter=0

while [ $counter -le $walletCnt ]
do
echo '1
1' | ../bin/ontology account add -d -w ../wallets/wallet$counter.dat > /dev/null
((counter++))
done

