build: build-client build-server

build-client:
	@echo starting creation client
	docker build --file=Dockerfile.client --tag=client:latest .

build-server:
	@echo starting creation server
	docker build --file=Dockerfile.server --tag=server:latest .

run: make-network run-server-d run-client

make-network:
	docker network create app || true
#Todo Open client and server in other system windows
run-client:
	@echo starting client
	docker run -it --name client --network=app client:latest --addr=server:7001

run-server:
	@echo starting server as daemon
	docker run --name server --network=app server:latest

run-server-d:
	@echo starting server as daemon
	docker run -d --name server --network=app server:latest

clean:
	@echo starting clear
	docker stop server client || true
	docker network rm app || true
	docker rm server client || true