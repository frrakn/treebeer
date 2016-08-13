#!/usr/bin/env bash

echo "========================================================" >> context.log
cd contextManager
go run manager.go >> ../context.log 2>&1 &
sleep 3
cd ../contextCache
go run cache.go >> ../context.log 2>&1 &
cd ../contextUpdater
go run updater.go >> ../context.log 2>&1 &
cd ../contextServer
go run server.go >> ../context.log 2>&1 &

