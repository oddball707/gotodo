package handler

import (
	"context"

	"github.com/sirupsen/logrus"

	pb "gotodo/proto"
	s "gotodo/service"
)

type Handler struct {
	logger  *logrus.Entry
	service s.Service
}

func NewHandler(logger *logrus.Entry, service s.Service) *Handler {
	return &Handler{
		logger:  logger.WithField("struct", "Handler"),
		service: service,
	}
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	h.logger.WithField("request", req).Trace("Create")

	result, err := h.service.Create(req.Task)
	if err != nil {
		h.logger.Error("failed to create todo item.", err)
		return nil, err
	}

	return &pb.CreateTaskResponse{Id: result.Id}, nil
}

func (h *Handler) Delete(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	h.logger.WithField("request", req).Trace("Delete")

	h.service.Delete(req.Id)

	return &pb.DeleteTaskResponse{}, nil
}

func (h *Handler) List(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	h.logger.Trace("Delete")

	tasks, err := h.service.Get()
	if err != nil {
		h.logger.Error("failed to get todo items.", err)
		return nil, err
	}

	return &pb.ListTasksResponse{Task: tasks}, nil
}

