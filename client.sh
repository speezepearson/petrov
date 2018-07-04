#!/bin/bash

ADDRESS=localhost:2344
NAME=$1

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

reader &

while true
do
    read
    echo "LAUNCH LAUNCH LAUNCH AAAAAAA"
    curl http://$ADDRESS/$NAME -d ''
done
