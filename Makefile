run:
	go run cmd/main.go

test:
	go test -v -count=1 ./...

coverage:
	go test -count=1 -coverprofile=tmp/coverage.out ./...
	go tool cover -func=tmp/coverage.out
