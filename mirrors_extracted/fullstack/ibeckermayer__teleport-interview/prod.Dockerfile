FROM node:12 AS frontend-builder
WORKDIR /app
COPY ./frontend /app/frontend
COPY ./certs /app/certs
RUN cd ./frontend && npm install
RUN cd ./frontend && npm run build

FROM golang:latest
WORKDIR /app
COPY --from=frontend-builder /app .
COPY ./backend /app/backend
RUN cd ./backend && go mod download
RUN cd ./backend && go build main.go
ENTRYPOINT cd ./backend/ && ./main