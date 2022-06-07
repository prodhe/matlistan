package main

import (
	"log"
	"net/http"

	"os"

	"github.com/prodhe/matlistan/handler"
	"github.com/prodhe/matlistan/template"
)

func main() {
	template.SetDirectory("./tmpl/")
	template.SetFuncMap(template.FuncMap{
		"halfway": func(i, j int) bool {
			if j%2 != 0 {
				j = j + 1
			}
			return i == j/2
		},
	})
	df := template.Fields{
		"StaticDir": "/static",
	}
	template.Load("index", df, "base.html", "index.html")
	template.Load("help", df, "base.html", "help.html")

	if os.Getenv("DEVMODE") != "" {
		log.Println("Development mode")
		template.Develop(true)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	h := handler.New()

	log.Println("Listening on port", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
