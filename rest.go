package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type submit struct {
	Code string
}

func hillRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hillName := ps.ByName("hill")

	hill, ok := hills[hillName]
	if !ok {
		w.Write([]byte("error"))
		return
	}
	w.Write(hill.JSON())
}

func hillSubmitRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hillName := ps.ByName("hill")

	log.Println(hillName, ": Posting Warrior")

	b, _ := ioutil.ReadAll(r.Body)
	s := submit{}
	json.Unmarshal(b, &s)

	warrior := parseWarrior(s.Code)

	if !warrior.Valid() {
		w.Write([]byte("Invalid Warrior. Name, Author, Type or Code may be missing!"))
		return
	}

	h, ok := hills[hillName]

	if !ok {
		w.Write([]byte("No or invalid Hill selected!"))
		return
	}

	if !h.Valid(warrior) {
		w.Write([]byte("Warrior too long or wrong type!"))
		return
	}

	if h.HasWarrior(warrior) {
		w.Write([]byte("Warrior already in Hill!"))
		return
	}

	log.Println(hillName, ":", warrior.Name, "Accepted")

	h.Queue <- warrior
	w.Write([]byte("Ok"))
}

func hillsRest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var names []string
	for key := range hills {
		names = append(names, key)
	}
	b, _ := json.Marshal(names)
	w.Write(b)
}

func initRest() {
	router := httprouter.New()

	router.GET("/hill/:hill", hillRest)
	router.POST("/hill/:hill", hillSubmitRest)
	router.GET("/hills", hillsRest)

	router.NotFound = http.FileServer(http.Dir("public"))

	log.Fatal(http.ListenAndServe(":8081", router))
}
