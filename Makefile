IMAGE_CHECK = $(shell docker image inspect --format="exists" go-starter-kit:latest)
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam,gomnd
GOIMPORTS_CMD=go run golang.org/x/tools/cmd/goimports
hello:
	echo "hello there!"

build:
	if [[ $(IMAGE_CHECK) == "exists" ]]; then \
		echo "Image Already Exists. Exitting" \
		exit 1; \
	elif [[ $(IMAGE_CHECK) != "exists" ]]; then \
		docker build --tag go-starter-kit .; \
	fi;

run: 
	docker run -d --env-file=.env -p 4000:4000 go-starter-kit

start: build run

tools:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint
	
lint:
	golangci-lint run --timeout 1m0s -v -E ${LINT_SETTINGS}

format:
	go mod tidy
	gofmt -s -w -l .
	${GOIMPORTS_CMD} -w .
