#!/bin/sh

rm .env
mv .env.prod .env

echo "I am here"
go build -o server

chmod +x ./server
./server