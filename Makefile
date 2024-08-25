.PHONY: start-dbdocker
start-docker:
	docker run -it -p 5432:5432 --name cubo-book -e POSTGRES_USER=cube -e POSTGRES_PASSWORD=secret postgres

.PHONY: start-docker
start-worker:
	docker run --rm golang-linux-build

.PHONY: start-server
start-server:
	CUBE_HOST=localhost CUBE_PORT=5555 go run main.go