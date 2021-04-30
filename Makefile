
all: docker-run

build: 
	go build -o ./src/ cmd/db/db.go 
	go build -o ./src/ cmd/api/api.go 
	go build -o ./src/ cmd/repo/repo.go 
	go build -o ./src/ cmd/web/web.go 

docker-build: 
	docker build -f images/api/Dockerfile -t bundle/api:1.0 . 
	docker build -f images/db/Dockerfile -t bundle/db:1.0 . 
	docker build -f images/web/Dockerfile -t bundle/web:1.0 . 
	docker build -f images/repo/Dockerfile -t bundle/repo:1.0 . 

docker-run: build docker-build
	docker-compose up -d


.PHONY: clean

clean:
	docker kill $(docker ps -q)