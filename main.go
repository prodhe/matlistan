package main

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"

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
	template.Load("signup", df, "base.html", "signup.html")
	template.Load("login", df, "base.html", "login.html")
	template.Load("help", df, "base.html", "help.html")
	template.Load("profile", df, "base.html", "profile.html")
	template.Load("recipes", df, "base.html", "recipes.html")

	if os.Getenv("DEVMODE") != "" {
		log.Println("Development mode")
		template.Develop(true)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "mongodb://localhost:27017/foodlist"
	}

	session, err := mgo.Dial(mongodbURI)
	if err != nil {
		log.Fatalf("Could not dial mongo db: %v\n", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	log.Println("Connected to MongoDB")

	h := handler.New(session.DB(""))

	log.Println("Listening on port", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
