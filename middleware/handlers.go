package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/mysticis/go-dcktst-demo/demo"
)

type Server struct {
	*http.ServeMux

	*demo.Queries
}

func NewServer() *Server {

	DB, err := initDB()

	if err != nil {
		log.Fatalln(err)
	}

	s := &Server{
		ServeMux: http.NewServeMux(),
		Queries:  DB,
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.HandleFunc("/create", s.CreateNewUser())
	s.HandleFunc("/update/", s.UpdateUser())
	s.HandleFunc("/getuser/", s.GetUser())
	s.HandleFunc("/getusers", s.GetAllUsers())
	s.HandleFunc("/delete/", s.DeleteUser())
	s.HandleFunc("/deleteusers", s.DeleteAllUsers())
}

//create
func (s *Server) CreateNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var u demo.User

		if r.Method != "POST" {

			http.Error(w, http.StatusText(405), 405)

			return
		}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		createdUser, err := s.Queries.CreateUser(context.Background(), demo.CreateUserParams{
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
		})

		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		json.NewEncoder(w).Encode(createdUser)
	}
}

//get a single user

func (s *Server) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return

		}
		id, err := strconv.Atoi(path.Base(r.URL.Path))

		if err != nil {
			return
		}

		retrievedUser, err := s.Queries.GetUser(context.Background(), int64(id))

		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(retrievedUser)
	}
}

//get all users

func (s *Server) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		storedUsers, err := s.Queries.ListUsers(context.Background())

		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(storedUsers)
	}
}

//update a task

func (s *Server) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var u demo.User

		if r.Method != "PUT" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		//read dynamic ID params
		id, err := strconv.Atoi(path.Base(r.URL.Path))
		if err != nil {
			return
		}

		//important
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		args := demo.UpdateUserParams{
			ID:    int64(id),
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
		}

		updatedUser, err := s.Queries.UpdateUser(context.Background(), args)

		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)

	}
}

//delete a user

func (s *Server) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		//get dynamic ID params
		id, err := strconv.Atoi(path.Base(r.URL.Path))
		if err != nil {
			return
		}

		err = s.Queries.DeleteUser(context.Background(), int64(id))

		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode("Deleted")
	}
}

//delete all records

func (s *Server) DeleteAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		err := s.Queries.DeleteAllUsers(context.Background())

		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return

		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode("All records deleted")
	}
}
