package ssl

import (
	"context"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ListSSLs = "list_ssls"
)

var ListSSLsTool = mcp.NewTool(
	ListSSLs,
	mcp.WithDescription("list ssls"),
)

func ListSSLHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	req := &types.PageRequest{
		Page:     1,
		PageSize: 500,
	}
	client := utils.NewPanelClient("POST", "/websites/ssl/search", utils.WithPayload(req))
	listWebsiteSSLRes := &types.ListWebsiteSSLRes{}
	return client.Request(listWebsiteSSLRes)
}
