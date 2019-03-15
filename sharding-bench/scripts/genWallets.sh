#!/bin/bash

walletCnt=20
counter=1

while [ $counter -le $walletCnt ]
do
echo '1
1' | ../bin/ontology account add -d -w ../wallets/wallet$counter.dat > /dev/null
((counter++))
done

