.PHONY: all vendor run

all: vendor run

vendor:
	go mod vendor

run:
	docker-compose -f deploy/docker-compose.yaml up --build
