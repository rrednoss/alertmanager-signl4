package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func NewServer() *http.Server {
	s := &http.Server{
		Addr:         ":9095",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      http.HandlerFunc(handleAlert),
	}
	return s
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "HEAD":
		handleHEAD(w, r)
	case "POST":
		handlePOST(w, r)
	default:
		http.Error(w, fmt.Sprintf("the HTTP method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func handleHEAD(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	if err := handlePOSTHeader(w, r); err != nil {
		return
	}
	if err := handlePOSTBody(w, r); err != nil {
		http.Error(w, fmt.Sprintf("error processing the request %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func handlePOSTHeader(w http.ResponseWriter, r *http.Request) error {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(w, "invalid Content-Type header", http.StatusUnsupportedMediaType)
		return errors.New("invalid Content-Type header")
	}
	return nil
}

func handlePOSTBody(w http.ResponseWriter, r *http.Request) error {
	alert, err := decodeBody(r.Body)
	if err != nil {
		return err
	}
	root := getRepositoryRootPath()
	gotpl, err := getTemplate(root, "signl4.gotpl")
	if err != nil {
		return err
	}
	tAlert, err := transform(gotpl, alert)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tAlert))

	return nil
}

func decodeBody(body io.ReadCloser) (map[string]interface{}, error) {
	var alert map[string]interface{}

	d := json.NewDecoder(body)
	if err := d.Decode(&alert); err != nil {
		return nil, err
	}
	return alert, nil
}

// Either I wasn't looking properly or it's pretty ridiculous. The Go libraries
// do not provide a reliable way to determine the root path of the repository.
// Using os.Getwd() the current path can be determined but it differs depending
// on whether I start the application or my test cases.
func getRepositoryRootPath() string {
	return os.Getenv("APP_ROOT_PATH")
}

func getTemplate(root string, name string) (string, error) {
	content, err := ioutil.ReadFile(root + "/templates/" + name)
	if err != nil {
		return "", fmt.Errorf(root+"/templates/"+name+" %w", err)
	}
	return string(content), nil
}
