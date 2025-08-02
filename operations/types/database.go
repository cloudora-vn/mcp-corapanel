package types

type Database struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`	
}

type DatabaseListResponse struct {
	Response
	Data struct {
		PageResult
		Items []Database `json:"items"`
	} `json:"data"`
} 

type ListDatabaseRequest struct {
	PageRequest
	Order    string `json:"order"`
	OrderBy  string `json:"orderBy"`
	Database string `json:"database"`
}


type CreateDatabaseRequest struct {
	Database string `json:"database"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Format   string `json:"format"`
	From     	string `json:"from"`
	Permission string `json:"permission"`
	Name 		string `json:"name"`
	Username 	string `json:"username"`
	Superuser 	bool `json:"superuser"`
}

