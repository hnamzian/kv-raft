package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

type HttpServer struct {
	addr string
	r  *raft.Raft
	db *sync.Map
}

func New(r *raft.Raft, db *sync.Map, addr string) *HttpServer {
	return &HttpServer{addr, r, db}
}

func (hs HttpServer) Start() {
	http.HandleFunc("/store", hs.keyhandler)
	http.HandleFunc("/join", hs.joinHandler)
	http.ListenAndServe(hs.addr, nil)
}

func (hs HttpServer) joinHandler(w http.ResponseWriter, r *http.Request) {
	followerId := r.URL.Query().Get("followerId")
	followerAddr := r.URL.Query().Get("followerAddr")

	if hs.r.State() != raft.Leader {
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			"Not the leader",
		})
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := hs.r.AddVoter(raft.ServerID(followerId), raft.ServerAddress(followerAddr), 0, 0).Error()
	if err != nil {
		log.Printf("Failed to add follower: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) keyhandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hs.setKey(w, r)
	case http.MethodGet:
		hs.getKey(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (hs HttpServer) setKey(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read key-value in http request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	future := hs.r.Apply(bs, 500*time.Millisecond)

	// Blocks until completion
	if err := future.Error(); err != nil {
		log.Printf("Could not write key-value: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	e := future.Response()
	if e != nil {
		log.Printf("Could not write key-value, application: %s", e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hs HttpServer) getKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, _ := hs.db.Load(key)
	if value == nil {
		value = ""
	}

	rsp := struct {
		Data string `json:"data"`
	}{value.(string)}
	err := json.NewEncoder(w).Encode(rsp)
	if err != nil {
		log.Printf("Could not encode key-value in http response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
