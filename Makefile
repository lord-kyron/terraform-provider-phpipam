test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/phpipam -run="TestAcc"

build: deps
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64 darwin/arm64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-phpipam" .

release: release_bump release_build

release_bump:
	scripts/release_bump.sh

release_build:
	scripts/release_build.sh

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
