package client

import (
	"net/http"
	"time"
)

type AlertStatus int

const (
	Firing AlertStatus = iota
	Resolved
)

type Signl4Client struct {
	Client     *http.Client
	FiringURL  string
	ResolveURL string
}

func (s Signl4Client) SendAlert(status AlertStatus) (int, error) {
	return 0, nil
}

func NewSignl4Client() Signl4Client {
	return Signl4Client{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		FiringURL:  "",
		ResolveURL: "",
	}
}
