#!/usr/bin/env bash

cd contextManager
go run manager.go >> ../context.log 2>&1 &
sleep 3
cd ../contextUpdater
go run updater.go >> ../context.log 2>&1 &
