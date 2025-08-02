package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var (
	accessToken string
	apiBase     string
	timestamp   string
)

func md5Sum(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func SetAccessToken(token string) {
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	accessToken = md5Sum("1panel" + token + timestamp)
}

func SetHost(host string) {
	apiBase = fmt.Sprintf("%s%s", host, ApiBase)
}

func GetAccessToken() string {
	if accessToken != "" {
		return accessToken
	}
	if token := os.Getenv("PANEL_ACCESS_TOKEN"); token != "" {
		SetAccessToken(token)
		return accessToken
	}
	return accessToken
}

func GetApiBase() string {
	if apiBase != "" {
		return apiBase
	}
	if host := os.Getenv("PANEL_HOST"); host != "" {
		SetHost(host)
		return apiBase
	}
	return apiBase
}

type PanelClient struct {
	Url       string
	Method    string
	Payload   interface{}
	Headers   map[string]string
	Response  *http.Response
	parsedUrl *url.URL
	Query     map[string]string
}

type Option func(client *PanelClient)

type ErrMsg struct {
	Message string `json:"message"`
}

type PanelError struct {
	Code    int
	Message string
	Details string
}

func (e *PanelError) Error() string {
	return fmt.Sprintf("Panel API error: %s (code: %d)", e.Message, e.Code)
}

func NewPanelError(code int, message, details string) *PanelError {
	return &PanelError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewAPIError(statusCode int, body []byte) error {
	var errMsg ErrMsg
	if err := json.Unmarshal(body, &errMsg); err != nil {
		details := string(body)
		if details == "" {
			details = "No error details available"
		}
		return NewPanelError(statusCode, http.StatusText(statusCode), details)
	}

	return NewPanelError(statusCode, http.StatusText(statusCode), errMsg.Message)
}

func NewAuthError() error {
	return NewPanelError(401, "Unauthorized", "Panel access token is missing or invalid")
}

func IsAuthError(err error) bool {
	var panelErr *PanelError
	if errors.As(err, &panelErr) {
		return panelErr.Code == 401
	}
	return false
}

func NewNetworkError(err error) error {
	return NewPanelError(0, "Network Error", err.Error())
}

func IsNetworkError(err error) bool {
	var panelErr *PanelError
	if errors.As(err, &panelErr) {
		return panelErr.Code == 0
	}
	return false
}

func NewInternalError(err error) error {
	return NewPanelError(500, "Internal Error", err.Error())
}

func IsAPIError(err error) bool {
	var panelError *PanelError
	ok := errors.As(err, &panelError)
	return ok
}

func NewPanelClient(method, urlPath string, opts ...Option) *PanelClient {
	urlString := GetApiBase() + urlPath
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	client := &PanelClient{
		Method:    method,
		Url:       parsedUrl.String(),
		parsedUrl: parsedUrl,
		Headers:   make(map[string]string),
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

func WithQuery(query map[string]interface{}) Option {
	return func(client *PanelClient) {
		parsedQuery := make(map[string]string)
		if query != nil {
			queryParams := client.parsedUrl.Query()
			for k, v := range query {
				parsedValue := ""
				switch v := v.(type) {
				case string:
					parsedValue = v
				case int:
					parsedValue = strconv.Itoa(v)
				case bool:
					parsedValue = strconv.FormatBool(v)
				}
				if parsedValue != "" {
					queryParams.Set(k, parsedValue)
					parsedQuery[k] = parsedValue
				}
			}
			client.parsedUrl.RawQuery = queryParams.Encode()
		}
		client.Url = client.parsedUrl.String()
		client.Query = parsedQuery
	}
}

func WithPayload(payload interface{}) Option {
	return func(client *PanelClient) {
		client.Payload = payload
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(client *PanelClient) {
		if client.Headers == nil {
			client.Headers = make(map[string]string)
		}
		for k, v := range headers {
			client.Headers[k] = v
		}
	}
}

func (p *PanelClient) SetHeaders(headers map[string]string) *PanelClient {
	if p.Headers == nil {
		p.Headers = make(map[string]string)
	}
	for k, v := range headers {
		p.Headers[k] = v
	}
	return p
}

func (p *PanelClient) Do() (*PanelClient, error) {
	p.Response = nil
	var reqBody io.Reader

	if p.Payload != nil {
		_payload, err := json.Marshal(p.Payload)
		if err != nil {
			return nil, NewInternalError(err)
		}
		reqBody = bytes.NewReader(_payload)
	}

	req, err := http.NewRequest(p.Method, p.Url, reqBody)
	if err != nil {
		return nil, NewInternalError(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "panel-client Go/"+runtime.GOOS+"/"+runtime.GOARCH+"/"+runtime.Version())

	token := GetAccessToken()
	if token == "" {
		return nil, NewAuthError()
	}

	req.Header.Set("1Panel-Token", token)
	req.Header.Set("1Panel-Timestamp", timestamp)

	for key, value := range p.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return p, NewNetworkError(err)
	}

	p.Response = resp

	if !p.IsSuccess() {
		body, _ := io.ReadAll(resp.Body)
		return p, NewAPIError(resp.StatusCode, body)
	}

	return p, nil
}

func (p *PanelClient) IsSuccess() bool {
	if p.Response == nil {
		return false
	}

	successMap := map[int]struct{}{
		http.StatusOK:          {},
		http.StatusCreated:     {},
		http.StatusNoContent:   {},
		http.StatusFound:       {},
		http.StatusNotModified: {},
	}

	_, ok := successMap[p.Response.StatusCode]
	return ok
}

func (p *PanelClient) IsFail() bool {
	return !p.IsSuccess()
}

func (p *PanelClient) GetRespBody() ([]byte, error) {
	if p.Response == nil || p.Response.Body == nil {
		return nil, errors.New("response or response body is nil")
	}
	defer p.Response.Body.Close()
	return io.ReadAll(p.Response.Body)
}

func (p *PanelClient) ParseJSON(v interface{}) error {
	body, err := p.GetRespBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (p *PanelClient) Request(object any) (*mcp.CallToolResult, error) {
	_, err := p.Do()
	if err != nil {
		switch {
		case IsAuthError(err):
			return mcp.NewToolResultText("Authentication failed: Please check your Panel access token"), err
		case IsNetworkError(err):
			return mcp.NewToolResultText("Network error: Unable to connect to Panel API"), err
		case IsAPIError(err):
			var panelErr *PanelError
			errors.As(err, &panelErr)
			return mcp.NewToolResultText(fmt.Sprintf("API error (%d): %s", panelErr.Code, panelErr.Details)), err
		default:
			return mcp.NewToolResultText(err.Error()), err
		}
	}

	if object == nil {
		return mcp.NewToolResultText("Operation completed successfully"), nil
	}

	body, err := p.GetRespBody()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Failed to read response body: %s", err.Error())),
			NewInternalError(err)
	}

	if err = json.Unmarshal(body, object); err != nil {
		errorMessage := fmt.Sprintf("Failed to parse response: %v", err)
		return mcp.NewToolResultText(errorMessage), NewInternalError(errors.New(errorMessage))
	}

	result, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Failed to format response: %s", err.Error())),
			NewInternalError(err)
	}

	return mcp.NewToolResultText(string(result)), nil
}
