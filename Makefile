
version = 1.0


all: docker-run

docker-build:
	go mod vendor
	docker-compose build 

docker-run: docker-build
	docker-compose up -d

dev:
	scripts/tls.sh
	scripts/run.sh


.PHONY: clean
clean:
	scripts/docker-kill.sh

