package types


type AppInstallCreate struct {
	AppDetailID uint                   `json:"appDetailId"`
	Params      map[string]interface{} `json:"params"`
	Name        string                 `json:"name"`
}

type AppRes struct {
	Response
	Data App `json:"data"`
}

type App struct {
	ID       uint   `json:"id"`
	Versions []string `json:"versions"`
}

type AppDetailRes struct {
	Response
	Data AppDetail `json:"data"`
}

type AppDetail struct {
	ID       uint   `json:"id"`
}
type AppInstall struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`	
	Status        string    `json:"status"`
	AppName       string    `json:"appName"`
}

type AppInstalledListResponse struct {
	Response
	Data struct {
		PageResult
		Items []AppInstall `json:"items"`
	} `json:"data"`
}