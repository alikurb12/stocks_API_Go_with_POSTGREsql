package main

import (
	"github.com/alikurb12/stocks_api_go/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080......")

	log.Fatal(http.ListenAndServe(":8008", r))
}