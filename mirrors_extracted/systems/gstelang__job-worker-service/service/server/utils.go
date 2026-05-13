package main

import (
	pb "github.com/gstelang/job-worker-service/service/proto/v1"
	"github.com/gstelang/job-worker-service/worker"
)

// Convert StartJobRequest to Command
func convertToCommand(req *pb.StartJobRequest) worker.Command {
	return worker.Command{
		Name: req.GetCommandName(),
		Args: req.GetCommandArgs(),
	}
}
