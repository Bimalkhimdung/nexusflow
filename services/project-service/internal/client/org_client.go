package client

import (
	"context"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/logger"
	orgv1 "github.com/nexusflow/nexusflow/pkg/proto/org/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrgClient wraps the org-service gRPC client
type OrgClient struct {
	client orgv1.OrgServiceClient
	conn   *grpc.ClientConn
	log    *logger.Logger
}

// NewOrgClient creates a new org-service client
func NewOrgClient(addr string, log *logger.Logger) (*OrgClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to org-service: %w", err)
	}

	client := orgv1.NewOrgServiceClient(conn)

	return &OrgClient{
		client: client,
		conn:   conn,
		log:    log,
	}, nil
}

// Close closes the connection
func (c *OrgClient) Close() error {
	return c.conn.Close()
}

// GetMemberRole gets a user's role in an organization
func (c *OrgClient) GetMemberRole(ctx context.Context, orgID, userID string) (orgv1.OrgRole, bool, error) {
	resp, err := c.client.GetMemberRole(ctx, &orgv1.GetMemberRoleRequest{
		OrganizationId: orgID,
		UserId:         userID,
	})
	if err != nil {
		c.log.Sugar().Errorw("Failed to get member role", "error", err, "org_id", orgID, "user_id", userID)
		return orgv1.OrgRole_ORG_ROLE_UNSPECIFIED, false, err
	}

	return resp.Role, resp.IsMember, nil
}

// IsAdmin checks if a user is an admin or owner in an organization
func (c *OrgClient) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	role, isMember, err := c.GetMemberRole(ctx, orgID, userID)
	if err != nil {
		return false, err
	}

	if !isMember {
		return false, nil
	}

	return role == orgv1.OrgRole_ORG_ROLE_ADMIN || role == orgv1.OrgRole_ORG_ROLE_OWNER, nil
}
