package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"slices"

	"github.com/xSudoNymx/job-manager/gen/jobpb"
	"github.com/xSudoNymx/job-manager/pkg/job"
	"github.com/xSudoNymx/job-manager/pkg/manager"
	"github.com/xSudoNymx/job-manager/pkg/output"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JobGRPCServer struct {
	jobpb.UnimplementedJobServiceServer
	mgr    *manager.Manager
	logger *slog.Logger
}

func NewJobGRPCServer(mgr *manager.Manager, logger *slog.Logger) (*JobGRPCServer, error) {
	if mgr == nil {
		return nil, errors.New("manager must not be nil")
	}

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return &JobGRPCServer{
		mgr:    mgr,
		logger: logger,
	}, nil
}

func (s *JobGRPCServer) Start(ctx context.Context, req *jobpb.StartRequest) (*jobpb.StartResponse, error) {
	if req.GetSpec() == nil {
		return nil, status.Error(codes.InvalidArgument, "missing process spec")
	}

	if req.GetSpec().GetExe() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing process path")
	}

	jobID, err := s.mgr.Start(req.GetSpec().GetExe(), slices.Clone(req.GetSpec().GetArgs()))
	if err != nil {
		return nil, toStatus(err)
	}

	return &jobpb.StartResponse{
		JobId: jobID,
	}, nil
}

func (s *JobGRPCServer) Stop(ctx context.Context, req *jobpb.StopRequest) (*jobpb.StopResponse, error) {
	if req.GetJobId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing job id")
	}

	if err := s.mgr.Stop(req.GetJobId()); err != nil {
		return nil, toStatus(err)
	}

	return &jobpb.StopResponse{}, nil
}

func (s *JobGRPCServer) Status(ctx context.Context, req *jobpb.StatusRequest) (*jobpb.StatusResponse, error) {
	if req.GetJobId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing job id")
	}

	st, err := s.mgr.Status(req.GetJobId())
	if err != nil {
		return nil, toStatus(err)
	}

	return &jobpb.StatusResponse{
		Status: toProtoJobStatus(st),
	}, nil
}

func (s *JobGRPCServer) StreamOutput(req *jobpb.StreamOutputRequest, stream jobpb.JobService_StreamOutputServer) error {
	if req.GetJobId() == "" {
		return status.Error(codes.InvalidArgument, "missing job id")
	}

	sub, err := s.mgr.Subscribe(req.GetJobId())
	if err != nil {
		return toStatus(err)
	}
	defer sub.Close()

	ctx := stream.Context()

	for {
		frame, err := sub.Read(ctx)
		switch {
		case err == nil:
			if sendErr := stream.Send(&jobpb.Output{
				Tag:  toProtoStreamTag(frame.Tag),
				Data: frame.Data,
			}); sendErr != nil {
				return sendErr
			}
			continue

		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil

		case errors.Is(err, io.EOF):
			return nil

		default:
			s.logger.Error("failed to read output frame", "error", err)
			return status.Error(codes.Internal, "failed to read output frame")
		}
	}
}

/* helpers */

func toProtoStreamTag(tag output.StreamTag) jobpb.StreamTag {
	switch tag {
	case output.TagStdout:
		return jobpb.StreamTag_STREAM_TAG_STDOUT
	case output.TagStderr:
		return jobpb.StreamTag_STREAM_TAG_STDERR
	default:
		return jobpb.StreamTag_STREAM_TAG_UNSPECIFIED
	}
}

func toProtoJobStatus(st *job.Status) *jobpb.JobStatus {
	return &jobpb.JobStatus{
		JobId:     st.ID,
		Pid:       int64(st.PID),
		State:     toProtoJobState(st.State),
		StartedAt: st.StartedAt.UnixNano(),
		EndedAt:   st.EndedAt.UnixNano(),
		ExitCode:  int32(st.ExitCode),
		WaitErr:   errToString(st.WaitErr),
	}
}

func toProtoJobState(s job.State) jobpb.JobState {
	switch s {
	case job.Running:
		return jobpb.JobState_JOB_STATE_RUNNING
	case job.Stopping:
		return jobpb.JobState_JOB_STATE_STOPPING
	case job.Stopped:
		return jobpb.JobState_JOB_STATE_STOPPED
	case job.Exited:
		return jobpb.JobState_JOB_STATE_EXITED
	default:
		return jobpb.JobState_JOB_STATE_UNSPECIFIED
	}
}

func errToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
