package app

import (
	"context"
	"fmt"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	InstallOpenResty = "install_openresty"
)

var InstallOpenRestyTool = mcp.NewTool(
	InstallOpenResty,
	mcp.WithDescription("install openresty, if not set name, default is openresty, if not set http_port, default is 80, if not set https_port, default is 443"),
	mcp.WithString("name", mcp.Description("openresty name"), mcp.DefaultString("openresty")),
	mcp.WithNumber("http_port", mcp.Description("openresty http port"), mcp.DefaultNumber(80)),
	mcp.WithNumber("https_port", mcp.Description("openresty https port"), mcp.DefaultNumber(443)),
)

func InstallOpenRestyHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var name string
	if request.Params.Arguments["name"] == nil {
		name = "openresty"
	} else {
		name = request.Params.Arguments["name"].(string)
	}

	var httpPort float64
	if request.Params.Arguments["http_port"] != nil {
		httpPort = request.Params.Arguments["http_port"].(float64)
	}
	if httpPort == 0 {
		httpPort = 80
	}

	var httpsPort float64
	if request.Params.Arguments["https_port"] != nil {
		httpsPort = request.Params.Arguments["https_port"].(float64)
	}
	if httpsPort == 0 {
		httpsPort = 443
	}

	appRes := &types.AppRes{}
	_, err := utils.NewPanelClient("GET", "/apps/openresty").Request(appRes)
	if err != nil {
		return nil, err
	}
	version := appRes.Data.Versions[0]
	appID := appRes.Data.ID
	appDetailUrl := fmt.Sprintf("/apps/detail/%d/%s/app", appID, version)
	appDetailRes := &types.AppDetailRes{}
	_, err = utils.NewPanelClient("GET", appDetailUrl).Request(appDetailRes)
	if err != nil {
		return nil, err
	}

	appDetailID := appDetailRes.Data.ID

	req := &types.AppInstallCreate{
		AppDetailID: appDetailID,
		Name:        name,
		Params: map[string]interface{}{
			"PANEL_APP_PORT_HTTP":  httpPort,
			"PANEL_APP_PORT_HTTPS": httpsPort,
		},
	}
	client := utils.NewPanelClient("POST", "/apps/install", utils.WithPayload(req))
	res := &types.Response{}
	return client.Request(res)
}
