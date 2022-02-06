package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type AlertStatus int

const (
	Firing AlertStatus = iota
	Resolved
	Unknown
)

type Signl4Client struct {
	Client     *http.Client
	FiringURL  string
	ResolveURL string
}

func (s Signl4Client) SendAlert(status AlertStatus, body io.Reader) (int, error) {
	req, err := http.NewRequest(http.MethodPost, s.getUrl(status), body)
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
	return res.StatusCode, nil
}

func (s Signl4Client) getUrl(status AlertStatus) string {
	if status == Firing {
		return s.FiringURL
	} else if status == Resolved {
		return s.ResolveURL
	}
	return ""
}

func NewSignl4Client() Signl4Client {
	// var tr = &http.Transport{}

	// if config.Signl4.AllowInsecureTLSConfig == "true" {
	// 	tr = &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	}
	// }

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return Signl4Client{
		Client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		},
		FiringURL:  "",
		ResolveURL: "",
	}
}
