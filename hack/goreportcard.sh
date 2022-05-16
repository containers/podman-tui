#!/usr/bin/env bash
set -e
git clone https://github.com/gojp/goreportcard.git
cd goreportcard
go install -mod=vendor ./vendor/github.com/alecthomas/gometalinter
go install -mod=vendor ./vendor/golang.org/x/lint/golint
go install -mod=vendor ./vendor/github.com/fzipp/gocyclo/cmd/gocyclo
go install -mod=vendor ./vendor/github.com/gordonklaus/ineffassign
go install -mod=vendor ./vendor/github.com/client9/misspell/cmd/misspell
go install ./cmd/goreportcard-cli
cd ../
rm -rf goreportcard
output=$($GOPATH/bin/goreportcard-cli | tee /dev/stdout)
echo $output | grep -q "A+"
if [ $? -ne 0 ] ; then
exit 1
fi
exit 0
