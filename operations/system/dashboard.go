package system

import (
	"context"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	GetDashboardInfo = "get_dashboard_info"
)

var GetDashboardInfoTool = mcp.NewTool(GetDashboardInfo, mcp.WithDescription(
	"show dashboard info"))

func GetDashboardInfoHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := utils.NewPanelClient("GET", "/dashboard/base/all/all")
	osInfo := &types.DashboardRes{}
	return client.Request(osInfo)
}
