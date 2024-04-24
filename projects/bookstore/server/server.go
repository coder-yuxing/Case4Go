package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"time"
	"yuxing.code/bookstore/server/middleware"

	"net/http"
	"yuxing.code/bookstore/store"
)

// BookStoreServer is a http server
type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, store store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		s: store,
		srv: &http.Server{
			Addr: addr,
		},
	}
	router := mux.NewRouter()
	router.HandleFunc("/book", srv.createBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", srv.updateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", srv.getBookHandler).Methods("GET")
	router.HandleFunc("/book", srv.getAllBooksHandler).Methods("GET")
	router.HandleFunc("/book/{id}", srv.delBookHandler).Methods("DELETE")
	srv.srv.Handler = middleware.Logging(middleware.Validating(router))
	return srv
}

func (s *BookStoreServer) createBookHandler(resp http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}

	if err := s.s.Create(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *BookStoreServer) updateBookHandler(resp http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(resp, "id is required", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	book.Id = id
	if err := s.s.Update(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *BookStoreServer) getBookHandler(resp http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(resp, "id is required", http.StatusBadRequest)
		return
	}

	book, err := s.s.Get(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	response(resp, book)
}

func (s *BookStoreServer) getAllBooksHandler(resp http.ResponseWriter, req *http.Request) {
	books, err := s.s.GetAll()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	response(resp, books)
}

func (s *BookStoreServer) delBookHandler(resp http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(resp, "id is required", http.StatusBadRequest)
		return
	}

	if err := s.s.Delete(id); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

// ListenAndServe starts the http server
func (s *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = s.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (s *BookStoreServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func response(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
