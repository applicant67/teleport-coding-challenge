# Build stage
FROM registry.access.redhat.com/ubi9/ubi:9.4-1214

# Install Go 1.23
RUN yum install -y wget && \
    wget https://go.dev/dl/go1.23.0.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.23.0.linux-arm64.tar.gz && \
    rm go1.23.0.linux-arm64.tar.gz

RUN yum install -y systemd procps-ng yum-utils && yum clean all

ENV PATH=$PATH:/usr/local/go/bin
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main service/server/main.go

# Removing multi-stage so I can exec and use golang although you do need a non-minimalistic linux distro
# such as Fedora to test memory or IO controller.

# Final stage
# FROM registry.access.redhat.com/ubi9/ubi:9.4-1214
# WORKDIR /app
# COPY --from=build /app/main .
# COPY certs/ certs/

EXPOSE 50051
CMD ["./main"]
