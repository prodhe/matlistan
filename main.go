package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"os"

	"github.com/prodhe/matlistan/handler"
	"github.com/prodhe/matlistan/template"
)

func main() {
	template.SetDirectory("./tmpl/")
	df := template.Fields{
		"StaticDir": "/static/",
	}
	template.Load("index", df, "base.html", "index.html")
	template.Load("signup", df, "base.html", "signup.html")
	template.Load("login", df, "base.html", "login.html")
	template.Load("about", df, "base.html", "about.html")
	template.Load("profile", df, "base.html", "profile.html")
	template.Load("recipes", df, "base.html", "recipes.html")

	template.Develop(true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		mongodb_uri = "mongodb://localhost:27017/foodlist"
	}

	session, err := mgo.Dial(mongodb_uri)
	if err != nil {
		fmt.Printf("could not dial mongo db: %v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	h := handler.New(session.DB(""))

	fmt.Println("Listening on port", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		fmt.Println(err)
	}
}
