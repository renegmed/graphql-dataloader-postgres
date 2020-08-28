init-project:
	go mod init github.com/tinrab/curly-waddle
    #go get github.com/99designs/gqlgen@v0.11.3
    # gqlgen version

init-gqlgen:
	go run github.com/99designs/gqlgen init
    # gqlgen init
    #gqlgen generate

generate:
	go generate ./...

up: 
	docker-compose up -d 

run:
	go run server.go 

