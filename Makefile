
compose = docker compose -p qonto -f docker/docker-compose.yml -f docker/docker-compose.dev.yml

help:
	@sed -n "/^[a-zA-Z0-9_-]*:/ s/:.*#/ -/p" < Makefile | sort

test: # Run short, non-integrational, tests
	@go test -race -short ./...

test-integration: # Run integration tests
	@go test -race ./...

test-all: test test-integration # Run all available tests, including integration

run: # Start project in background
	${compose} up -d database

mysql-enter: # Run mysql client inside database container
	${compose} exec database mysql -uroot -proot

logs: # Follow logs from all containers in the project
	${compose} logs -f

down: # Stop all project containers
	${compose} stop

clean: # Remove project containers
	${compose} rm -fs
