package ssl

import (
	"context"
	"errors"
	"github.com/cloudora-vn/mcp-corapanel/operations/types"
	"github.com/cloudora-vn/mcp-corapanel/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	CreateSSL = "create_ssl"
)

var CreateSSLTool = mcp.NewTool(
	CreateSSL,
	mcp.WithDescription("create ssl"),
	mcp.WithString("domain", mcp.Description("domain"), mcp.Required()),
	mcp.WithString("provider", mcp.Description("provider support dnsAccount,http"), mcp.Required(), mcp.DefaultString("dnsAccount")),
	mcp.WithString("dnsAccount", mcp.Description("dnsAccount"), mcp.DefaultString("")),
)

func CreateSSLHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if request.Params.Arguments["domain"] == nil {
		return nil, errors.New("domain is required")
	}
	domain := request.Params.Arguments["domain"].(string)
	if request.Params.Arguments["provider"] == nil {
		return nil, errors.New("provider is required")
	}
	provider := request.Params.Arguments["provider"].(string)
	if provider != "dnsAccount" && provider != "http" {
		return nil, errors.New("provider must be dnsAccount or http")
	}

	acmeRes := &types.ListAcmeRes{}
	pageReq := &types.PageRequest{
		Page:     1,
		PageSize: 500,
	}
	client := utils.NewPanelClient("POST", "/websites/acme/search", utils.WithPayload(pageReq))
	_, err := client.Request(acmeRes)
	if err != nil {
		return nil, err
	}
	if len(acmeRes.Data.Items) == 0 {
		return nil, errors.New("no acme account found")
	}
	acme := acmeRes.Data.Items[0]

	var dnsAccountID uint
	if provider == "dnsAccount" {
		dnsAccountRes := &types.ListDNSAccountRes{}
		client = utils.NewPanelClient("POST", "/websites/dns/search", utils.WithPayload(pageReq))
		_, err = client.Request(dnsAccountRes)
		if err != nil {
			return nil, err
		}
		if len(dnsAccountRes.Data.Items) == 0 {
			return nil, errors.New("no dns account found")
		}
		var dnsName string
		if request.Params.Arguments["dnsAccount"] != nil {
			dnsName = request.Params.Arguments["dnsAccount"].(string)
		}
		if dnsName != "" {
			checkName := strings.ToLower(dnsName)
			for _, dnsAccount := range dnsAccountRes.Data.Items {
				if strings.Contains(strings.ToLower(dnsAccount.Name), checkName) || strings.Contains(strings.ToLower(dnsAccount.Type), checkName) {
					dnsAccountID = dnsAccount.ID
					break
				}
			}
		}
		if dnsAccountID == 0 {
			dnsAccountID = dnsAccountRes.Data.Items[0].ID
		}
	}

	req := &types.CreateSSLRequest{
		PrimaryDomain: domain,
		Provider:      provider,
		AcmeAccountID: acme.ID,
		DnsAccountID:  dnsAccountID,
		KeyType:       "P256",
	}
	client = utils.NewPanelClient("POST", "/websites/ssl", utils.WithPayload(req))
	res := &types.Response{}
	return client.Request(res)
}
