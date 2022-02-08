package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"webserver/internal/app/model"
	"webserver/internal/app/store"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName               = "ebweb"
	contextKeyUser contextKey = iota
	contextKeyRequestID
)

var (
	errorIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errorNotAuthenticated         = errors.New("not authenticated")
)

type contextKey int8

type server struct {
	router       *mux.Router
	store        store.Store
	logger       *logrus.Logger
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
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
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()

	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		rw.Header().Set("X-Request-ID", id)
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), contextKeyRequestID, id)))
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)

		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]

		if !ok {
			s.error(rw, r, http.StatusUnauthorized, errorNotAuthenticated)
		}

		u, err := s.store.User().Find(id.(int))

		if err != nil {
			s.error(rw, r, http.StatusUnauthorized, errorNotAuthenticated)
			return
		}

		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), contextKeyUser, u)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		loger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(contextKeyRequestID),
		})
		loger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		customRW := &responseWriter{rw, http.StatusOK}

		next.ServeHTTP(customRW, r)

		loger.Infof("completed with %d %s in %v",
		customRW.code,
		http.StatusText(customRW.code),
		time.Since(start))
	})
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.respond(rw, r, http.StatusOK, r.Context().Value(contextKeyUser).(*model.User))
	}
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
