package system

import (
	"context"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	GetSystemInfo = "get_system_info"
)

var GetSystemInfoTool = mcp.NewTool(GetSystemInfo, mcp.WithDescription(
	"show host system information, The unit of diskSize is bytes"))

func GetSystemInfoHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := utils.NewPanelClient("GET", "/dashboard/base/os")
	osInfo := &types.OsInfoRes{}
	return client.Request(osInfo)
}
