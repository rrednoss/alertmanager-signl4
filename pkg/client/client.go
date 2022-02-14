package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	SendAlert(status AlertStatus, body io.Reader) (int, error)
}

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

func NewSignl4Client() Signl4Client {
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

func (sc Signl4Client) SendAlert(status AlertStatus, body io.Reader) (int, error) {
	req, err := http.NewRequest(http.MethodPost, sc.getUrl(status), body)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	fmt.Println("Send alert to ", sc.getUrl(status))
	fmt.Println(req)
	res, err := sc.Client.Do(req)
	if err != nil {
		return 0, err
	}
	if res.StatusCode == http.StatusOK {
		return res.StatusCode, nil
	}
	return res.StatusCode, nil
}

func (sc Signl4Client) getUrl(status AlertStatus) string {
	if status == Firing {
		return sc.FiringURL
	} else if status == Resolved {
		return sc.ResolveURL
	}
	return ""
}
