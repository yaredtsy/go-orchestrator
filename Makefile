.PHONY start-db
start-db:
	docker run -it -p 5432:5432 --name cubo-book -e POSTGRES_USER=cube -e POSTGRES_PASSWORD=secret postgres

.PHONY start-docker
start-worker:
	docker run --rm golang-linux-build