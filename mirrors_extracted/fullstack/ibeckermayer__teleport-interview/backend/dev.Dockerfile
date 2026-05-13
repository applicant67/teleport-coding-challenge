FROM golang:latest

WORKDIR /app

COPY ./backend /app
COPY ./certs /certs

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix="false" --build="go build main.go" --command="./main --env=dev"