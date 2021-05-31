
version = 1.0

all: docker-run

docker-build: tidy sass cli
	docker-compose build 

docker-run: docker-build
	docker-compose up -d

dev: tidy clean sass cli
	./dev.sh

sass:
	sass ./assets/public/scss/styles.scss ./assets/public/css/styles.css

cli:
	go install bundle.go

tidy:
	go mod tidy
	go mod vendor

.ONESHELL: dev clean cert

cert:
	./cert.sh


proto:
	protoc --gofast_out=plugins=grpc:. ./api/api.proto


.PHONY: clean cli
clean:
	docker-compose down
	./clean.sh

full-clean:
	./clean.sh
	docker system prune -a --volumes


