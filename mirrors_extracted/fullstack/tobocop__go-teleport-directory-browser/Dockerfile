# Build the Go API
FROM arm64v8/golang:latest AS builder
ADD . /app
WORKDIR /app/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# Build the React application
FROM arm64v8/node:14.17.5-alpine3.11 AS node_builder
COPY --from=builder /app/client ./
RUN yarn install
RUN yarn run build
# Final stage build, this will be the container
# that we will deploy to production
FROM arm64v8/alpine:latest
COPY --from=builder /main ./
COPY --from=builder /app/server/certs ./certs
COPY --from=node_builder /build ./web
RUN chmod +x ./main
EXPOSE 8080
CMD ./main
