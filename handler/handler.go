package handler

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/prodhe/matlistan/model"
	"github.com/prodhe/matlistan/template"
)

type handler struct {
	mux *http.ServeMux
	db  *mgo.Database
}

func New(db *mgo.Database) *handler {
	h := &handler{
		mux: http.NewServeMux(),
		db:  db,
	}

	h.mux.HandleFunc("/signup", h.signup)
	h.mux.HandleFunc("/recipes", h.recipes)

	h.mux.HandleFunc("/", h.index)

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) recipes(w http.ResponseWriter, r *http.Request) {
	ingredients := make([]string, 0)
	ingredients = append(ingredients, "10 liter vatten")
	recipe := model.Recipe{
		Title:       "Grönsakssoppa",
		Category:    "Vegetariskt",
		Ingredients: ingredients,
		Description: "Koka vatten och drick.",
	}
	fmt.Fprintf(w, "%v", recipe)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./web/")).ServeHTTP(w, r)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "signup", nil)
}
