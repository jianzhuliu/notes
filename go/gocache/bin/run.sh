#!/bin/bash

##用于在 shell 脚本退出时，删掉临时文件，结束子进程
trap "rm ./bin/multi-node;kill 0" EXIT

go build cmd/multi-node.go -o bin/
./bin/multi-node -port=8001 &
./bin/multi-node -port=8002 &
./bin/multi-node -port=8003 &
./bin/multi-node -api &

sleep 2
echo ">>> start test"
curl "http://127.0.0.1:9999/api?key=k4" &
curl "http://127.0.0.1:9999/api?key=k7" &
curl "http://127.0.0.1:9999/api?key=k6" &

wait