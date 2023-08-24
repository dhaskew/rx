detected_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')

#.SILENT: lint build

all: vet lint static build

setup:
	go get ./...

up:
	docker compose up -d

down:
	docker compose down

vet:
	go vet ./...

complexity:
	#functions higher than 9 are an issue
	~/go/bin/gocyclo -avg -ignore "_test|Godeps|vendor/" .

list_files:
	#show the files we will include in the build
	go list -f '{{ .GoFiles }}' ./...

lint:
	~/go/bin/golangci-lint run

static:
	staticcheck ./internal/*

build:
	go build -o rx -race ./main.go

test:
	go test -race -v -shuffle=on -coverprofile cover.out ./...

#doc:
#	open http://localhost:8080/github.com/dhaskew/rx &
#	pkgsite

coverage: test
	go tool cover -html=cover.out -o cover.html
	go tool cover -func=cover.out
ifeq ($(detected_OS), Linux)
	xdg-open cover.html
endif
ifeq ($(detected_OS), Darwin)
	open cover.html
endif

clean:
	rm -f rx
	rm -f cover.out
	rm -f cover.html