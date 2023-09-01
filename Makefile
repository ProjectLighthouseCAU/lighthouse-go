run:
	go run main.go

build:
	go build -o lighthouse-go

build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -installsuffix cgo -o lighthouse-go-arm

full-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lighthouse-go
