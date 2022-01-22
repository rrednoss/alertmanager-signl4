package main

import "github.com/rrednoss/alertmanager-signl4/pkg/server"

func main() {
	s := server.NewServer()
	s.ListenAndServe()
}
