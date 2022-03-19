package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"simple_app_for_kube/cmd/database"
)

type handler struct {
	db database.DBManager
}

func NewHandler(db *sql.DB) *handler {
	return &handler{db: database.NewDBManager(db)}
}

func (h *handler) configureRouter(r *mux.Router) {
	r.HandleFunc("/users", h.handleGetUserName()).Methods("GET")
	r.HandleFunc("/users", h.handleUserCreate()).Methods("POST")
}

func (h *handler) handleUserCreate() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req = request{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.error(w, http.StatusBadRequest, err)
			return
		}

		var (
			userID int64
			err    error
		)
		if userID, err = h.db.UserRepo().Save(req.Name); err != nil {
			h.error(w, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, http.StatusCreated, userID)
	}
}

func (h *handler) handleGetUserName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			userIDStr = r.URL.Query().Get("userId")
			userID    int
			err       error
		)

		if userID, err = strconv.Atoi(userIDStr); err != nil {
			h.error(w, http.StatusBadRequest, err)
			return
		}

		var name string
		if name, err = h.db.UserRepo().GetName(int64(userID)); err != nil {
			h.error(w, http.StatusInternalServerError, err)
			return
		}
		h.respond(w, http.StatusOK, name)
	}
}

func (h *handler) error(w http.ResponseWriter, code int, err error) {
	h.respond(w, code, map[string]string{"error": err.Error()})
	log.Println("error happened:", err)
}

func (h *handler) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println("can't encode data", err)
		}
	}
}
