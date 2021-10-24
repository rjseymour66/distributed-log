package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Accepts an IP address string and returns a pointer to an HTTP 
// server
func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()

	// Handles POST and GET requests to the '/' home endpoint
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:		addr,
		Handler:	r,
	}
}

type httpServer struct {
	Log *Log
}

// Server that references a Log for the server
func newHTTPServer() *httpServer {
	return &httpServer{
		Log:	NewLog(),
	}
}

// Contains the record that the user wants to append to the log
// in the POST request
type ProduceRequest struct {
	Record Record `json:"record"`
}

// The offset returned to a successful POST request
type ProduceResponse struct {
	Offset uint64	`json:"offset"`
}

// The offset that represents a record that the user wants
// to retrieve in a GET request
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// The record returned to a successful GET request
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
