package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func BackpackIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusBadRequest, Text: "No username provided"}); err != nil {
		panic(err)
	}
}

func BackpackShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	var limit int
	var err error

	if username, ok := vars["username"]; ok {
		// set limit to 1 if it is not set
		if limit, err = strconv.Atoi(r.FormValue("limit")); err != nil {
			limit = 1
			fmt.Println("Limit must be an integer")
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		backpack_array := RepoFindBackpack(username, limit)
		if len(backpack_array) > 0 {
			w.WriteHeader(http.StatusOK)
			var hostname string
			if host, err := os.Hostname(); err != nil {
				hostname = "unknown"
				panic(err)
			} else {
				hostname = host
			}
			if err := json.NewEncoder(w).Encode(ReturnVal{backpack_array, hostname}); err != nil {
				panic(err)
			}
		} else {
			// If we didn't find it, 404
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}

type ReturnVal struct {
	Backpacks Backpacks `json:"backpacks"`
	Hostname  string    `json:"hostname"`
}

func BackpackAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	backpack := Backpack{0, r.FormValue("username"), "", r.FormValue("backpack_json")}

	b := RepoAddBackpack(backpack)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		panic(err)
	}
}
