
version = 1.0


all: docker-run


docker-build:
	go mod vendor
	docker-compose build

docker-run: docker-build
	docker-compose up -d


.PHONY: clean
clean:
	docker-compose down

