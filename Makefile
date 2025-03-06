.PHONY: run test coverage api_test docker-build docker-up docker-down

run:
	go run cmd/main.go

test:
	go test -v -count=1 ./...

coverage:
	go test -count=1 -coverprofile=tmp/coverage.out ./...
	go tool cover -func=tmp/coverage.out

api_test:
	hurl --error-format=long --verbose --test --variable host=localhost:8080 api.hurl

# Docker commands
docker-down:
	docker-compose down

# Combined commands
docker:
	docker-compose build app
	docker-compose up app

# Run tests in Docker
docker-test:
	docker-compose run --rm app go test -v ./...

# Run coverage in Docker
docker-coverage:
	docker-compose run --rm app go test ./... -coverprofile=tmp/coverage.out
	docker-compose run --rm app go tool cover -func=tmp/coverage.out

# Run API tests in Docker
docker-api-test:
	docker-compose run --rm hurl --error-format=long --verbose --test --variable host=app:8080 /app/api.hurl
