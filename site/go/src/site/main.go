package main

import (
	"net/http"
	router "site/app/route"
)

func main() {
	router := router.GetRouter()
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:80",
	}

	server.ListenAndServe()
}
