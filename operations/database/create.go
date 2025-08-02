package database

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	CreateDatabase = "create_database"
)

var CreateDatabaseTool = mcp.NewTool(
	CreateDatabase,
	mcp.WithDescription("create a database by type name and password"),
	mcp.WithString("database_type", mcp.Description("installed database app type, support mysql and postgresql"), mcp.DefaultString("mysql"), mcp.Required()),
	mcp.WithString("database", mcp.Description("installed database app name"), mcp.DefaultString(""), mcp.Required()),
	mcp.WithString("name", mcp.Description("database name"), mcp.DefaultString(""), mcp.Required()),
	mcp.WithString("username", mcp.Description("database username"), mcp.DefaultString("")),
	mcp.WithString("password", mcp.Description("database password"), mcp.DefaultString("")),
)

func CreateDatabaseHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var (
		database     string
		password     string
		name         string
		databaseType string
		username     string
	)
	if request.Params.Arguments["database"] == nil {
		return nil, errors.New("database name is required")
	}
	database = request.Params.Arguments["database"].(string)
	if request.Params.Arguments["database_type"] == nil {
		return nil, errors.New("database type is required")
	}
	databaseType = request.Params.Arguments["database_type"].(string)
	if databaseType != "mysql" && databaseType != "postgresql" {
		return nil, errors.New("database type is invalid, support mysql and postgresql")
	}
	if request.Params.Arguments["name"] == nil {
		return nil, errors.New("name is required")
	}
	name = request.Params.Arguments["name"].(string)
	if request.Params.Arguments["password"] == nil {
		password = utils.GetRandomStr(12)
	} else {
		password = request.Params.Arguments["password"].(string)
	}
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	if request.Params.Arguments["username"] == nil {
		username = name
	} else {
		username = request.Params.Arguments["username"].(string)
	}

	createReq := &types.CreateDatabaseRequest{
		Database: database,
		Password: encodedPassword,
		Type:     databaseType,
		Name:     name,
		From:     "local",
		Username: username,
	}
	var createUrl string
	if databaseType == "mysql" {
		createUrl = "/databases"
		createReq.Format = "utf8mb4"
		createReq.Permission = "%"
	} else {
		createUrl = "/databases/pg"
		createReq.Format = "UTF8"
	}
	client := utils.NewPanelClient("POST", createUrl, utils.WithPayload(createReq))
	res := &types.Response{}
	return client.Request(res)
}
