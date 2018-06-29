SRCS = frame.go main.go reader.go tui.go utils.go
VERSION_FLAGS=-ldflags "-X main.version=`cat VERSION` -X main.date=`date -u +%Y/%m/%d-%H:%M:%S`"
ARGS ?= -h ; printf "\n***\n* Set ARGS variable in make invocation\n***"

all: build

build:
	go build $(VERSION_FLAGS) $(SRCS)

clean:
	go clean

run:
	go run $(VERSION_FLAGS) $(SRCS) $(ARGS)

test:
	go test -v

update:
	dep ensure -v -update

format:
	find . -path ./vendor -prune -o -name '*.go' -exec gofmt -s -w {} \;
