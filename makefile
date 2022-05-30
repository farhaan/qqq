prepare:
	go mod tidy;
	go test ./...;
	go build ./main.go