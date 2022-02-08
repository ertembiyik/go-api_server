package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"webserver/internal/app/model"
	"webserver/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "ebweb"
)
var (
	errorIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionCreate()).Methods("POST")
}

func (s *server) handleUserCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(rw, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		s.respond(rw, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(rw, r, http.StatusUnauthorized, errorIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)

		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
		}

		session.Values["user_id"] = u.ID

		if err := s.sessionStore.Save(r, rw, session); err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
		}
		
		s.respond(rw, r, http.StatusOK, nil)
	}
}

func (s *server) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(rw, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)

	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}
