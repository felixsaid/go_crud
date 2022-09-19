package main

import (
	"fmt"
	"log"
	"net/http"
	"v1/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on port 8000.........")

	log.Fatal(http.ListenAndServe(":8000", r))
}
