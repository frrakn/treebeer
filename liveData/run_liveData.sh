#!/usr/bin/env bash

echo "========================================================" >> liveData.log
cd translator
echo "Starting Translator"
go run translator.go >> ../liveData.log 2>&1 &
echo "Waiting 20 seconds then starting Aggregator..."
sleep 20
cd ../aggregator
echo "Starting Aggregator..."
go run aggregator.go >> ../liveData.log 2>&1 &

