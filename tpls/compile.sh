#!/usr/bin/env sh

go-bindata -o tpls.go microservice/... client/...

sed -i "" "s/package main/package tpls/g" tpls.go
