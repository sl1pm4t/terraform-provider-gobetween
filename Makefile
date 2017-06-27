test:
	go get ./...
	go get github.com/dustinkirkland/golang-petname
	go test -timeout 20m -v ./lxd

testacc:
	TF_LOG=debug TF_ACC=1 go test ./gobetween -v $(TESTARGS)

build:
	go build -v
	tar czvf terraform-provider-gobetween_${TRAVIS_TAG}_linux_amd64.tar.gz terraform-provider-gobetween

dev:
	go build -v
