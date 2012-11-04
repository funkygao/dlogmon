#!/bin/sh
echo ========================
echo stdlib benchmark
echo ========================
go test -test.bench=".*" fmt
go test -test.bench=".*" bytes
go test -test.bench=".*" strconv
go test -test.bench=".*" regexp

echo
echo

for f in `ls benchm*.go`
do
    echo $f
    echo ========================
    go run $f
    echo
done
