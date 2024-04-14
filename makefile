VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X main.Version=v-$(VERSION)
GOLDFLAGS += -X main.BuildTime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"


linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o ./bin/linux-monitor-amd64 $(GOFLAGS) ./main/

linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build  -o ./bin/linux-monitor-arm64 $(GOFLAGS) .

darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build  -o ./bin/darwin-monitor-arm64 $(GOFLAGS) .

darwin-amd64: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -buildvcs=false  -o ./bin/darwin-monitor-amd64 $(GOFLAGS) ./main/

windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -o ./bin/windows-monitor-amd64 $(GOFLAGS) .

clean:
	rm -f bin/*monitor*
	rm -fr ./bin/logs/*.log
all: linux-amd64 linux-arm64 darwin-arm64 darwin-amd64 windows-amd64
	
