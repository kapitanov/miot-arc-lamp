package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func runHttp() {
	r := mux.NewRouter()

	r.HandleFunc("/state", httpGetState).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www")))

	go func() {
		http.ListenAndServe(":3000", r)
	}()
}

func httpGetState(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(currentStatus)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
