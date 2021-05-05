
version = 1.0

define dev-script = 
	export DATABASE_PORT=8040
	export REPO_PORT=8060
	export WEB_PORT=8080
	export GATE_PORT=8020
	export GATE_HOST=localhost
	export WEB_HOST=localhost
	export REPO_HOST=localhost
	export DATABASE_HOST=localhost
	export MODE=DEV
	export AUTH0_ID=MpOxXFrk5XhR7gKcWIhYVZTNDDinx4ZT
	export AUTH0_SECRET=jdksX1I0hZ8vej4M6LW-VRxtIiRFVXr2MMVYK0K9FvD8EtsiiRfATnKszcb2SvrG
	export AUTH0_AUD=https://bundlemc.io/auth
	export MONGO_URL="mongodb+srv://benny-bundle:thisismypassword1@bundle.mveuj.mongodb.net/main?retryWrites=true&w=majority"
	export MONGO_AUTH=FALSE
	export AWS_REGION=us-east-1
	export AWS_BUCKET=bundle-repository
	go run cmd/gate/gate.go &
	go run cmd/db/db.go &
	go run cmd/repo/repo.go & 
	go run cmd/web/web.go && fg
endef

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
	$(value dev-script)


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


