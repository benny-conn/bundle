
version = 1.0



define clean-script =
	docker-compose down
	docker rm -f $(docker ps -a -q)
	docker volume rm $(docker volume ls -q)
	fuser -k 8020/tcp 
	fuser -k 8040/tcp 
	fuser -k 8060/tcp 
	fuser -k 8080/tcp 
endef




all: docker-run

docker-build:
	go mod vendor
	docker-compose build 

docker-run: docker-build
	docker-compose up -d

dev:
	./dev.sh


.ONESHELL: dev clean cert

cert:
	./cert.sh


proto:
	protoc --gofast_out=plugins=grpc:. ./api/api.proto


.PHONY: clean
clean:
	$(value clean-script)

full-clean:
	$(value clean-script)
	docker system prune -a --volumes


