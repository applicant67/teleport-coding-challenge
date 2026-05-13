proto:
	protoc -I=. --go_out=./internal/worker --go-grpc_out=./internal/worker ./proto/worker.proto

test:
	go test -race -v -timeout 30s -failfast -cover github.com/jaylane/job-scheduler/...