package types

import "time"

type ListWebsiteSSLRes struct {
	Response
	Data struct {
		PageResult
		Items []WebsiteSSL `json:"items"`
	} `json:"data"`
}

type WebsiteSSL struct {
	ID            uint      `json:"id"`
	PrimaryDomain string    `json:"primaryDomain"`
	Domains       string    `json:"domains"`
	Provider      string    `json:"provider"`
	Organization  string    `json:"organization"`
	AutoRenew     bool      `json:"autoRenew"`
	ExpireDate    time.Time `json:"expireDate"`
	StartDate     time.Time `json:"startDate"`
	Status        string    `json:"status"`
}

type CreateSSLRequest struct {
	PrimaryDomain string `json:"primaryDomain"`
	Domains       string `json:"domains"`
	Provider      string `json:"provider"`
	AcmeAccountID uint   `json:"acmeAccountId"`
	DnsAccountID  uint   `json:"dnsAccountId"`
	KeyType       string `json:"keyType"`
}

type ListAcmeRes struct {
	Response
	Data  AcmeDTO `json:"data"`
}

type AcmeDTO struct {
	PageResult
	Items []Acme `json:"items"`
}

type Acme struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

type ListDNSAccountRes struct {
	Response
	Data struct {
		PageResult
		Items []DNSAccount `json:"items"`
	} `json:"data"`
}

type DNSAccount struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}