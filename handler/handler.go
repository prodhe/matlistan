package handler

import (
	"fmt"
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
	h.mux.HandleFunc("/session", h.session)

	h.mux.HandleFunc("/", h.index)

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	template.Render(w, "index", nil)
}

func (h *handler) session(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionGet(w, r)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	fmt.Fprintf(w, "%v", session)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		template.Render(w, "signup", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only GET/POST allowed.", http.StatusMethodNotAllowed)
	}

	fmt.Printf("%v\n", r)

	fmt.Printf("username: %v\npassword: %v\n", r.PostFormValue("username"), r.PostFormValue("password"))

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

func (h *handler) sessionGet(w http.ResponseWriter, r *http.Request) (*model.Session, error) {
	c, err := r.Cookie("sid")
	fmt.Printf("cookie get: %v\n", c)
	if err != nil || c.Value == "" {
		return h.sessionSet(w, r)
	}

	sid, err := url.QueryUnescape(c.Value)
	if err != nil || !bson.IsObjectIdHex(sid) {
		return nil, fmt.Errorf("could not unescape session id: %v", err)
	}

	count, err := h.db.C("session").FindId(bson.ObjectIdHex(sid)).Count()
	if count < 1 || err != nil {
		return h.sessionSet(w, r)
	}

	var session model.Session

	err = h.db.C("session").FindId(bson.ObjectIdHex(sid)).One(&session)
	if err != nil {
		return nil, fmt.Errorf("could not find session: %v", err)
	}

	return &session, nil
}

func (h *handler) sessionSet(w http.ResponseWriter, r *http.Request) (*model.Session, error) {
	session := model.Session{
		Id:            bson.NewObjectId(),
		Pid:           bson.NewObjectId(),
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
	fmt.Printf("cookie set: %v\n", cookie)
	http.SetCookie(w, &cookie)

	return &session, nil
}
