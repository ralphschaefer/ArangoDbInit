version = 0.2.0
name = arangoinit
local_tag = $(name):$(version)
remote_tag = ralphschaefer/arangodb-cloudinit:$(version)

all: build

docker:
	docker build -t $(local_tag) .

docker-publish: docker
	docker tag $(local_tag) $(remote_tag)
	docker push $(remote_tag)

build:
	rm -f $(name)
	CGO_ENABLED=0 GOOS=linux go build -o $(name) -ldflags="-s -w" main.go

run:
	go generate
	go run main.go

clean:
	rm -f main $(name) *~
