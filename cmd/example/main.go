package main

import (
	"log"
	"net/http"

	https "github.com/pav5000/golang-https"
)

func main() {
	srv := https.New("./data", "my@email.com", "my.site.com")

	err := srv.ListenHTTPS(":443", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	}))
	if err != nil {
		log.Fatal(err)
	}
}
