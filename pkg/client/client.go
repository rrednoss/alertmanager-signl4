package client

import (
	"io"
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

func (s Signl4Client) SendAlert(status AlertStatus, body io.Reader) (int, error) {
	req, err := http.NewRequest(http.MethodPost, s.FiringURL, body)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := s.Client.Do(req)
	if err != nil {
		return 0, err
	}
	if res.StatusCode == http.StatusOK {
		return res.StatusCode, nil
	}
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
