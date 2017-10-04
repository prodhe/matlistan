package handler

import (
	"context"
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

	h.mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))

	h.mux.HandleFunc("/signup", h.signup)
	h.mux.HandleFunc("/login", h.login)
	h.mux.HandleFunc("/logout", h.logout)
	h.mux.HandleFunc("/session", h.session)

	h.mux.HandleFunc("/about", h.about)

	h.mux.HandleFunc("/deleteaccount", h.sessionHandle(h.deleteAccount))
	h.mux.HandleFunc("/profile", h.sessionHandle(h.profile))
	h.mux.HandleFunc("/recipes/add", h.sessionHandle(h.recipesAdd))
	h.mux.HandleFunc("/recipes", h.sessionHandle(h.recipes))

	h.mux.HandleFunc("/", h.sessionHandle(h.index))

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

	h.db.C("recipe").Insert(recipe)

	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}

func (h *handler) about(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "about", nil)
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

func (h *handler) sessionHandle(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
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
