#!/bin/bash

ADDRESS=$1
NAME=$2

trap 'kill $READER_PID; exit' SIGINT SIGTERM

reader() {
    while true
    do
        clear
        date
        echo $NAME
        curl http://$ADDRESS/$NAME
        sleep 1
    done
}

reader & READER_PID=$!

while true
do
    read
    echo "LAUNCH LAUNCH LAUNCH AAAAAAA"
    curl http://$ADDRESS/$NAME -d ''
done
