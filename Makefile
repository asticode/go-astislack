build:
	GOOS=darwin go build -o dist/darwin cmd/main.go
	GOOS=linux go build -o dist/linux cmd/main.go
	GOOS=windows go build -o dist/windows cmd/main.go