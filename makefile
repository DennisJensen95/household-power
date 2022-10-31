
test:
	go test -coverprofile=coverage.out ./...

build: 
	go build -o bin/household-power cmd/power_usage/main.go
