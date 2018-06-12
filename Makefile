NAME = exchanger

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)

LD_OPTS = -ldflags="-X main.branch=${BRANCH} -X main.commit=${COMMIT} -X main.buildtime=${BUILDTIME} -w"

all:  build run

all-with-deps: setup deps build

run: build
	cd cmd && ./$(NAME) && ../

setup:
	go get -u github.com/kardianos/govendor

deps:
	govendor sync

protobuf:
	cd stream/server/ && protoc --go_out=plugins=grpc:. *.proto

build:
	cd ./cmd/ && go build $(LD_OPTS) -o $(NAME)  . && cd -

race:
	cd ./cmd/ && go build $(LD_OPTS) -o $(NAME) -race . && cd -

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow|nolint:.*' .
.PHONY: todo

dist:
	cd cmd/ && GOOS=linux GOARCH=amd64 go build $(LD_OPTS)  -o $(NAME) .