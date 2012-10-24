#!/bin/sh
for f in `ls benchm*.go`
do
    echo $f
    echo ========================
    go run $f
    echo
done
