package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prodhe/matlistan/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		mongodb_uri = "mongodb://localhost:27017/foodlist"
	}

	h := handler.New()

	fmt.Println("Listening on port", port)
	err := http.ListenAndServe(":"+port, h)
	if err != nil {
		fmt.Println(err)
	}
}
