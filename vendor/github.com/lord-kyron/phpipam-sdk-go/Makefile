.PHONY: test testacc

test: deps
	go test -v ./...

testacc: deps
	TESTACC=1 go test -p 1 -v ./... -run="TestAcc"

deps:
	go get -u github.com/kardianos/govendor
	govendor sync
