# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o job-worker

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/job-worker /app/
COPY ssl /app/ssl
ENTRYPOINT ./job-worker server
