package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/haydenjeune/kvstore/pkg/store"
)

type GetRequest struct {
	Key string `json:"key"`
}

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func makeGetEndpointFunc(store store.KvStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body GetRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		value, exists, err := store.Get(body.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if !exists {
			http.NotFound(w, r)
			return
		}

		response := GetResponse{Key: body.Key, Value: value}
		encoder := json.NewEncoder(w)
		encoder.Encode(response)
	}
}

func makeSetEndpointFunc(store store.KvStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SetRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = store.Set(body.Key, body.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// TODO: Make the storage engine configurable
	//store, err := store.NewInMemHashMapKVStorage()
	//store, err := store.NewFsAppendOnlyStorage("data.kvstore")
	store, err := store.NewHashIndexedFsAppendOnlyStorage("data.kvstore")
	//store, err := store.NewInMemSortedKVStorage()
	if err != nil {
		log.Fatalf("Failed to instantiate storage: %v", err)
	}

	http.HandleFunc("/get", makeGetEndpointFunc(store))
	http.HandleFunc("/set", makeSetEndpointFunc(store))

	addr := "127.0.0.1:8080"
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
