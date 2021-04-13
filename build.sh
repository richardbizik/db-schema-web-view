#!/bin/sh

cd _webapp
npm install
npm run build
cd ..

go mod download
go mod verify
go build -o db-schema-web-view ./main