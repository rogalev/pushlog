build:
	go mod download && CGO_ENABLED=0 go build -o ./.bin/pushlog ./cmd/app/main.go