#!/usr/bin/env bash

cd contextManager
go run manager.go >> ../context.log 2>&1 &
#go run manager.go | awk '{ print strftime("%Y-%m-%dT%H:%M:%S"), $0; fflush(); }' >> ../output.log 2>&1 &
sleep 3
cd ../contextUpdater
go run updater.go >> ../context.log 2>&1 &
