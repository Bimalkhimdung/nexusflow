package handler

import (
	"github.com/nexusflow/nexusflow/pkg/logger"
	// Import your generated protobuf package here
)

// Handler implements the gRPC service
type Handler struct {
	// Embed UnimplementedServer for forward compatibility
	// pb.UnimplementedYourServiceServer
	
	service Service
	log     *logger.Logger
}

// Service interface defines business logic operations
type Service interface {
	// Define your service methods here
}

// NewHandler creates a new gRPC handler
func NewHandler(service Service, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// Example RPC method implementation
// func (h *Handler) GetExample(ctx context.Context, req *pb.GetExampleRequest) (*pb.GetExampleResponse, error) {
// 	h.log.Info("GetExample called", logger.String("id", req.Id))
// 	
// 	// Call service layer
// 	result, err := h.service.GetExample(ctx, req.Id)
// 	if err != nil {
// 		h.log.Error("Failed to get example", logger.Error(err))
// 		return nil, status.Error(codes.Internal, "failed to get example")
// 	}
// 	
// 	return &pb.GetExampleResponse{
// 		Example: result,
// 	}, nil
// }
