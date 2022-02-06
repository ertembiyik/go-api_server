package apiserver

import (
	"encoding/json"
	"io"
	"net/http"
	"webserver/internal/app/model"
	"webserver/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting apiserver")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/users", s.getUsers())
	s.router.HandleFunc("/create_user", s.create_note())
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)

	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) getUsers() http.HandlerFunc {
	// here you can define request specific types, variables etc

	return func(rw http.ResponseWriter, r *http.Request) {

		io.WriteString(rw, "Hello")

		// notes, err := s.store.User().GetAll()

		// if err != nil {
		// 	rw.Header().Set("Content-Type", "application/json")
		// 	io.WriteString(rw, "server error")
		// }

		// rw.Header().Set("Content-Type", "application/json")
		// json.Marshal(notes)
		// json.NewEncoder(rw).Encode(notes)
	}
}

func (s *APIServer) create_note() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		note := model.User{Email: "ertembiyik@gmail.com", Password: "Body"}
		noter, err := s.store.User().Create(&note)

		if err != nil {
			rw.Header().Set("Content-Type", "application/json")
			io.WriteString(rw, "server error")
		}

		rw.Header().Set("Content-Type", "application/json")
		json.Marshal(noter)
		json.NewEncoder(rw).Encode(noter)
	}
}
