package handler

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/attachment/v1"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/models"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AttachmentHandler struct {
	pb.UnimplementedAttachmentServiceServer
	svc *service.AttachmentService
	log *logger.Logger
}

func NewAttachmentHandler(svc *service.AttachmentService, log *logger.Logger) *AttachmentHandler {
	return &AttachmentHandler{svc: svc, log: log}
}

// UploadAttachment handles streaming file upload
func (h *AttachmentHandler) UploadAttachment(stream pb.AttachmentService_UploadAttachmentServer) error {
	var metadata *models.UploadMetadata
	var buffer bytes.Buffer

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			h.log.Sugar().Errorw("Failed to receive chunk", "error", err)
			return status.Errorf(codes.Internal, "failed to receive chunk: %v", err)
		}

		switch data := req.Data.(type) {
		case *pb.UploadAttachmentRequest_Metadata:
			metadata = &models.UploadMetadata{
				EntityType:  data.Metadata.EntityType,
				EntityID:    data.Metadata.EntityId,
				Filename:    data.Metadata.Filename,
				ContentType: data.Metadata.ContentType,
				UploaderID:  data.Metadata.UploaderId,
			}
		case *pb.UploadAttachmentRequest_Chunk:
			if _, err := buffer.Write(data.Chunk); err != nil {
				h.log.Sugar().Errorw("Failed to write chunk", "error", err)
				return status.Errorf(codes.Internal, "failed to write chunk: %v", err)
			}
		}
	}

	if metadata == nil {
		return status.Error(codes.InvalidArgument, "metadata is required")
	}

	// Validate file
	maxSize := int64(50 * 1024 * 1024) // 50MB
	if err := service.ValidateFile(int64(buffer.Len()), metadata.ContentType, maxSize, nil); err != nil {
		return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Upload file
	attachment, err := h.svc.UploadAttachment(stream.Context(), metadata, buffer.Bytes())
	if err != nil {
		h.log.Sugar().Errorw("Failed to upload attachment", "error", err)
		return status.Errorf(codes.Internal, "failed to upload attachment: %v", err)
	}

	return stream.SendAndClose(&pb.UploadAttachmentResponse{
		Attachment: attachmentToProto(attachment),
	})
}

// GetAttachment retrieves attachment metadata
func (h *AttachmentHandler) GetAttachment(ctx context.Context, req *pb.GetAttachmentRequest) (*pb.GetAttachmentResponse, error) {
	attachment, err := h.svc.GetAttachment(ctx, req.Id)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get attachment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get attachment: %v", err)
	}

	return &pb.GetAttachmentResponse{Attachment: attachmentToProto(attachment)}, nil
}

// GetDownloadURL generates a presigned download URL
func (h *AttachmentHandler) GetDownloadURL(ctx context.Context, req *pb.GetDownloadURLRequest) (*pb.GetDownloadURLResponse, error) {
	url, expiresIn, err := h.svc.GetDownloadURL(ctx, req.AttachmentId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get download URL", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get download URL: %v", err)
	}

	return &pb.GetDownloadURLResponse{
		Url:       url,
		ExpiresIn: expiresIn,
	}, nil
}

// ListAttachments lists attachments for an entity
func (h *AttachmentHandler) ListAttachments(ctx context.Context, req *pb.ListAttachmentsRequest) (*pb.ListAttachmentsResponse, error) {
	attachments, err := h.svc.ListAttachments(ctx, req.EntityType, req.EntityId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list attachments", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list attachments: %v", err)
	}

	var pbAttachments []*pb.Attachment
	for _, a := range attachments {
		pbAttachments = append(pbAttachments, attachmentToProto(a))
	}

	return &pb.ListAttachmentsResponse{Attachments: pbAttachments}, nil
}

// DeleteAttachment deletes an attachment
func (h *AttachmentHandler) DeleteAttachment(ctx context.Context, req *pb.DeleteAttachmentRequest) (*pb.DeleteAttachmentResponse, error) {
	if err := h.svc.DeleteAttachment(ctx, req.AttachmentId); err != nil {
		h.log.Sugar().Errorw("Failed to delete attachment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete attachment: %v", err)
	}

	return &pb.DeleteAttachmentResponse{}, nil
}

// Helper function
func attachmentToProto(a *models.Attachment) *pb.Attachment {
	if a == nil {
		return nil
	}
	return &pb.Attachment{
		Id:               a.ID,
		EntityType:       a.EntityType,
		EntityId:         a.EntityID,
		Filename:         a.Filename,
		OriginalFilename: a.OriginalFilename,
		ContentType:      a.ContentType,
		Size:             a.Size,
		UploaderId:       a.UploaderID,
		CreatedAt:        a.CreatedAt.Format(time.RFC3339),
	}
}
