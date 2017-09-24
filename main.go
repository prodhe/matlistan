package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Listening on port", port)
	err := http.ListenAndServe(":"+port, http.FileServer(http.Dir("./web/")))
	if err != nil {
		fmt.Println(err)
	}
}
