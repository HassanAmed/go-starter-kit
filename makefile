IMAGE_CHECK = $(shell docker image inspect --format="exists" go-starter-kit:latest)
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
	go get github.com/kisielk/errcheck
	go get github.com/golang/lint/golint
	go get github.com/axw/gocov/gocov
	go get github.com/matm/gocov-html
	go get github.com/tools/godep
	go get github.com/mitchellh/gox