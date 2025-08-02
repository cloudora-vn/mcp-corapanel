package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	InstallMySQL = "install_mysql"
)

var InstallMySQLTool = mcp.NewTool(
	InstallMySQL,
	mcp.WithDescription("install mysql, if not set name, default is mysql, if not set version, default is '', if not set root_password, default is '')"),
	mcp.WithString("name", mcp.Description("mysql name")),
	mcp.WithString("version", mcp.Description("mysql version, not support latest version"), mcp.DefaultString("")),
	mcp.WithString("root_password", mcp.Description("mysql root password"), mcp.DefaultString("")),
	mcp.WithNumber("port", mcp.Description("mysql port"), mcp.DefaultNumber(3306)),
)

func InstallMySQLHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var name string
	if request.Params.Arguments["name"] == nil {
		name = "mysql"
	} else {
		name = request.Params.Arguments["name"].(string)
	}
	var version string
	if request.Params.Arguments["version"] != nil {
		version = request.Params.Arguments["version"].(string)
		if version == "latest" {
			version = ""
		}
	}
	appRes := &types.AppRes{}
	_, err := utils.NewPanelClient("GET", "/apps/mysql").Request(appRes)
	if err != nil {
		return nil, err
	}
	exist := false
	for _, v := range appRes.Data.Versions {
		if v == version || strings.Contains(v, version) {
			version = v
			exist = true
			break
		}
	}
	if !exist {
		return nil, errors.New("version not found")
	}
	if version == "" {
		version = appRes.Data.Versions[0]
	}
	appID := appRes.Data.ID
	appDetailUrl := fmt.Sprintf("/apps/detail/%d/%s/app", appID, version)
	appDetailRes := &types.AppDetailRes{}
	_, err = utils.NewPanelClient("GET", appDetailUrl).Request(appDetailRes)
	if err != nil {
		return nil, err
	}
	appDetailID := appDetailRes.Data.ID
	var port float64
	if request.Params.Arguments["port"] != nil {
		port = request.Params.Arguments["port"].(float64)
	}
	if port == 0 {
		port = 3306
	}
	var rootPassword string
	if request.Params.Arguments["root_password"] != nil {
		rootPassword = request.Params.Arguments["root_password"].(string)
	}
	if rootPassword == "" {
		rootPassword = fmt.Sprintf("mysql_%s", utils.GetRandomStr(6))
	}

	req := &types.AppInstallCreate{
		AppDetailID: appDetailID,
		Name:        name,
		Params: map[string]interface{}{
			"PANEL_APP_PORT_HTTP":    port,
			"PANEL_DB_ROOT_PASSWORD": rootPassword,
		},
	}
	client := utils.NewPanelClient("POST", "/apps/install", utils.WithPayload(req))
	res := &types.Response{}
	return client.Request(res)
}
