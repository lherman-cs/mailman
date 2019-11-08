#!/bin/bash

endpoint="ws://localhost:4001/ws/"
n=10000

for i in `seq $n`
do
    websocat -H="Name: $i" $endpoint > /dev/null
done

sleep 10000
wait