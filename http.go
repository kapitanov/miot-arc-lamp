package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	httpEndpoint = "0.0.0.0:3000"
)

func runHttp() {
	r := mux.NewRouter()

	r.HandleFunc("/api/state", httpGetState).Methods("GET")
	r.HandleFunc("/api/state", httpSetState).Methods("POST")
	r.HandleFunc("/api/state", httpResetState).Methods("DELETE")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www")))

	go func() {
		fmt.Fprintf(os.Stdout, "http: listening on \"%s\"\n", httpEndpoint)
		http.ListenAndServe(httpEndpoint, r)
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

func httpSetState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var body ArcStatus
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		panic(err)
	}

	setTestStatus(&body)
	mqttPublish(&body)
	w.WriteHeader(204)
}

func httpResetState(w http.ResponseWriter, r *http.Request) {
	resetStatus()

	go func() {
		reqUpdate <- 0
	}()
	w.WriteHeader(204)
}
