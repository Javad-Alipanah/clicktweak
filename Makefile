APP=bin/clicktweak
STATIC_APP=bin/clicktweak_static
SOURCES=$(shell find . -name '*.go' -not -name '*_test.go')

all: $(APP)

static: $(STATIC_APP)

format:
	find . -name '*.go' -not -path "./.cache/*" | xargs -n1 go fmt

check: format
	git diff
	git diff-index --quiet HEAD

lint:
	golangci-lint run --skip-dirs=test --deadline 3m0s

test:
	go test -cover ./... -coverprofile .coverage.txt
	cat .coverage.txt | grep "/pkg\|mode:" > .coverage.pkg
	go tool cover -func .coverage.pkg

clean:
	rm -rf bin

bin/%: cmd/% $(SOURCES)
	go build -o $@ $</*.go
	strip -s $@

bin/%_static: cmd/% $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ $</*.go
	strip -s $@

bin:
	mkdir -p $@

$(TEST_DIR):
	mkdir -p $@

.PHONY: all static format check lint test clean
