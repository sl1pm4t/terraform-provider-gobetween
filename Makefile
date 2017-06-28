test:
	go get ./...
	go get github.com/dustinkirkland/golang-petname
	go test -timeout 20m -v ./gobetween

testacc:
	$(eval DOCKER_ID := $(shell docker run --rm -d -p 8888:8888 -v `pwd`/test/:/etc/gobetween/conf/:rw yyyar/gobetween))
	TF_LOG=info TF_ACC=1 GB_HOST=localhost GB_PORT=8888 go test ./gobetween -v $(TESTARGS)
	docker stop ${DOCKER_ID}
	

build:
	go build -v
	tar czvf terraform-provider-gobetween_${TRAVIS_TAG}_linux_amd64.tar.gz terraform-provider-gobetween

dev:
	go build -v

install:
	go install -v
