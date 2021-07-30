all: docker-build

CWD := $(shell pwd)
POSTGRES_DOCKER_IMAGE_BASE ?= postgresql-fast
PG_DOCKER_TAG ?= 13-alpine

POSTGRES_DOCKER_IMAGE := $(POSTGRES_DOCKER_IMAGE_BASE):$(PG_DOCKER_TAG)

docker-build:
	docker build -t postgresql-fast:13-alpine .

database-up:
	@echo Using $(POSTGRES_DOCKER_IMAGE)
	@docker run \
		--name boundary-sql-tests \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=boundary \
		-e POSTGRES_USER=boundary \
		-e POSTGRES_DB=boundary \
		-e PGDATA=/pgdata \
		--mount type=tmpfs,destination=/pgdata \
		-v "$(CWD)/../schema/migrations/postgres":/migrations \
		$(POSTGRES_DOCKER_IMAGE)
clean:
	docker stop boundary-sql-tests || true
	docker rm -v boundary-sql-tests || true

.PHONY: all docker-build database-up clean
