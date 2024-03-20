package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

var noteStore = make(map[string]Note)
var id int

func main() {
	m := mux.NewRouter().StrictSlash(false)
	m.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	m.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	m.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	m.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")
	server := http.Server{
		Addr:           ":8080",
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	e := server.ListenAndServe()
	if e != nil {
		log.Fatal(e)
	}
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Creating a Slice of Note struct.
	var notes []Note
	for _, v := range noteStore {
		notes = append(notes, v)
	}
	// Creating HTTP Header.
	w.Header().Set("Content-Type", "application/json")
	// Marshal noteStore values.
	sb, e := json.Marshal(notes)
	if e != nil {
		log.Fatal(e)
	}
	// Writing HTTP Header.
	w.WriteHeader(http.StatusOK)
	_, e = w.Write(sb)
	if e != nil {
		log.Fatal(e)
	}
}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Creating a new Note struct
	var note Note
	// Decoding the HTTP Body from Request
	e := json.NewDecoder(r.Body).Decode(&note)
	if e != nil {
		log.Fatal(e)
	}
	//	Setting Creation Date to Note
	note.CreatedAt = time.Now()
	// Creating New ID to Note
	id++
	// Converting ID to Str
	k := strconv.Itoa(id)
	// Adding note to noteStore in the position of ID
	noteStore[k] = note

	//	Returning the Created Object
	w.Header().Set("Content-Type", "application/json")
	sb, e := json.Marshal(note)
	if e != nil {
		log.Fatal(e)
	}
	w.WriteHeader(http.StatusCreated)
	_, e = w.Write(sb)
	if e != nil {
		log.Fatal(e)
	}
}

func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Convert the Request Variables to a map[string]string
	vars := mux.Vars(r)
	// k Stores the Note struct of the id Value
	k := vars["id"]
	// Creating Note struct
	var noteUpdate Note
	// Decoding the HTTP Body from Request
	e := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if e != nil {
		log.Fatal(e)
	}
	// Checking if the Note noteUpdate with ID k exists
	if note, ok := noteStore[k]; ok {
		noteUpdate.CreatedAt = note.CreatedAt
		delete(noteStore, k)
		noteStore[k] = noteUpdate
	} else {
		log.Printf("The Note with ID %s don't exists", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Convert the Request Variables to a map[string]string
	vars := mux.Vars(r)
	// k Stores the Note struct of the id Value
	k := vars["id"]
	// Checking if the Note noteUpdate with ID k exists
	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
	} else {
		log.Printf("The Note with ID %s don't exists", k)
	}
	w.WriteHeader(http.StatusNoContent)
}
