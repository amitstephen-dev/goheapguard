.PHONY: test bench clean build

test:
go test -race -cover ./...

bench:
go test -bench=. -benchmem ./...

build:
go build ./...

clean:
go clean -cache
rm -f coverage.out coverage.html

coverage:
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
