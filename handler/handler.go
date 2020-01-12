package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

	h.mux.Handle("/static/", h.noCache(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/")))))

	h.mux.HandleFunc("/signup", h.signup)
	h.mux.HandleFunc("/login", h.login)
	h.mux.HandleFunc("/logout", h.logout)
	h.mux.HandleFunc("/session", h.session)
	h.mux.HandleFunc("/help", h.help)

	h.mux.HandleFunc("/deleteaccount", h.sessionValidate(h.deleteAccount))
	h.mux.HandleFunc("/profile", h.sessionValidate(h.profile))
	h.mux.HandleFunc("/recipes/add", h.sessionValidate(h.recipesAdd))
	h.mux.HandleFunc("/recipes/delete", h.sessionValidate(h.recipesDelete))
	h.mux.HandleFunc("/recipes", h.sessionValidate(h.recipes))
	h.mux.HandleFunc("/recipes.json", h.sessionValidate(h.apiRecipes))

	h.mux.HandleFunc("/", h.sessionValidate(h.index))

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	data := template.Fields{
		"Authenticated": true,
	}
	template.Render(w, "index", data)
}

func (h *handler) help(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionGet(w, r)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	data := template.Fields{
		"Authenticated": session.Authenticated,
	}
	template.Render(w, "help", data)
}

func (h *handler) profile(w http.ResponseWriter, r *http.Request) {
	data := template.Fields{
		"Authenticated": true,
	}
	template.Render(w, "profile", data)
}

func (h *handler) recipes(w http.ResponseWriter, r *http.Request) {
	pid := r.Context().Value("pid")

	recipes := make([]model.Recipe, 0)
	err := h.db.C("recipe").Find(bson.M{"pid": pid}).All(&recipes)
	if err != nil {
		recipes = nil
	}

	data := template.Fields{
		"Authenticated": true,
		"Recipes":       recipes,
	}
	template.Render(w, "recipes", data)
}

func (h *handler) apiRecipes(w http.ResponseWriter, r *http.Request) {
	pid := r.Context().Value("pid")

	recipes := make([]model.Recipe, 0)
	err := h.db.C("recipe").Find(bson.M{"pid": pid}).All(&recipes)
	if err != nil {
		recipes = nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}

func (h *handler) recipesAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed.", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.sessionGet(w, r)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	rid := r.PostFormValue("rid")
	formtitle := r.PostFormValue("title")
	formingredients := r.PostFormValue("ingredients")
	description := r.PostFormValue("description")

	recipe := model.Recipe{
		Id:          bson.NewObjectId(),
		Pid:         session.Pid,
		Description: description,
	}
	title, categories := recipe.BreakTitle(formtitle)
	ingredients := recipe.BreakIngredients(formingredients)
	recipe.Title = title
	recipe.Categories = categories
	recipe.Ingredients = ingredients

	if rid == "" {
		h.db.C("recipe").Insert(recipe)
		log.Println("inserted new recipe:", recipe.Title)
	} else {
		recipe.Id = bson.ObjectIdHex(rid)
		err := h.db.C("recipe").Update(bson.M{"_id": recipe.Id}, recipe)
		if err != nil {
			log.Println("could not update recipe:", err)
		}
	}

	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}

func (h *handler) recipesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed.", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.sessionGet(w, r)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	rid := r.PostFormValue("rid")

	if rid == "" {
		http.Error(w, "ID empty.", http.StatusBadRequest)
		return
	}

	id := bson.ObjectIdHex(rid)

	recipes := make([]model.Recipe, 0)
	if err := h.db.C("recipe").Find(bson.M{"_id": id}).All(&recipes); err != nil {
		http.Error(w, "Not found.", http.StatusNotFound)
		return
	}
	rec := recipes[0]

	if rec.Pid != session.Pid {
		http.Error(w, "Not found.", http.StatusNotFound)
		return
	}

	if err := h.db.C("recipe").RemoveId(id); err != nil {
		log.Println("could not remove recipe:", err)
	} else {
		log.Println("deleted:", rec.Title)
	}

	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		template.Render(w, "signup", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only GET/POST allowed.", http.StatusMethodNotAllowed)
		return
	}

	u := r.PostFormValue("username")
	p, err := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), 10)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}
	pid := bson.NewObjectId()

	count, err := h.db.C("account").Find(bson.M{"username": u}).Count()
	if count > 0 || err != nil {
		data := template.Fields{
			"FormError": "Användarnamnet finns redan.",
			"Username":  u,
		}
		template.Render(w, "signup", data)
		return
	}

	account := model.Account{
		Id:       bson.NewObjectId(),
		Pid:      pid,
		Username: u,
		Password: string(p),
	}
	profile := model.Profile{
		Id:   pid,
		Name: u,
	}
	h.db.C("account").Insert(account)
	h.db.C("profile").Insert(profile)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		template.Render(w, "login", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only GET/POST allowed.", http.StatusMethodNotAllowed)
	}

	u := r.PostFormValue("username")
	p := r.PostFormValue("password")

	var account model.Account

	err := h.db.C("account").Find(bson.M{"username": u}).One(&account)
	err2 := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(p))
	if err != nil || err2 != nil {
		data := template.Fields{
			"FormError": "Fel användarnamn/lösenord.",
			"Username":  u,
		}
		template.Render(w, "login", data)
		return
	}

	session, err := h.sessionGet(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("session error: %s", err), http.StatusInternalServerError)
	}
	err = h.db.C("session").UpdateId(session.Id, bson.M{
		"pid":           account.Pid,
		"lastseen":      time.Now(),
		"authenticated": true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("session error: %s", err), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	h.sessionDelete(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// FIXME
func (h *handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed.", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) session(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionGet(w, r)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
	}

	fmt.Fprintf(w, "session: %v", session)
}

func (h *handler) sessionValidate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionGet(w, r)
		if err != nil {
			log.Fatal("session error")
		}

		if !session.Authenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "pid", session.Pid)

		next(w, r.WithContext(ctx))
	}
}

func (h *handler) sessionGet(w http.ResponseWriter, r *http.Request) (*model.Session, error) {
	var session model.Session

	c, err := r.Cookie("sid")
	if err != nil || c.Value == "" {
		return h.sessionSet(w, r)
	}

	sid, err := url.QueryUnescape(c.Value)
	if err != nil || !bson.IsObjectIdHex(sid) {
		return nil, fmt.Errorf("could not unescape session id: %v", err)
	}

	dbsess := h.db.Session.Copy()
	defer dbsess.Close()

	count, err := dbsess.DB("").C("session").FindId(bson.ObjectIdHex(sid)).Count()
	if count < 1 || err != nil {
		return h.sessionSet(w, r)
	}

	err = dbsess.DB("").C("session").FindId(bson.ObjectIdHex(sid)).One(&session)
	if err != nil {
		return nil, fmt.Errorf("could not find session: %v", err)
	}

	return &session, nil
}

func (h *handler) sessionSet(w http.ResponseWriter, r *http.Request) (*model.Session, error) {
	session := model.Session{
		Id:            bson.NewObjectId(),
		Pid:           "000000000000",
		LastSeen:      time.Now(),
		Authenticated: false,
	}

	err := h.db.C("session").Insert(&session)
	if err != nil {
		return nil, fmt.Errorf("could not create session: %v", err)
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    url.QueryEscape(session.Id.Hex()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   0,
	}
	http.SetCookie(w, &cookie)

	return &session, nil
}

func (h *handler) sessionDelete(w http.ResponseWriter, r *http.Request) error {
	session, err := h.sessionGet(w, r)
	if err != nil {
		return fmt.Errorf("could not get session: %v", err)
	}

	err = h.db.C("session").RemoveId(session.Id)
	if err != nil {
		return fmt.Errorf("could not remove session id: %v", err)
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   0,
	}
	http.SetCookie(w, &cookie)

	return nil
}

func (h *handler) noCache(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		next.ServeHTTP(w, r)
	}
}
