package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("starting...")
	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64 `json:"id"`
	PlaceID int64 `json:"place_id"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(w http.ResponseWriter, r *http.Request) {
	chk := Check{}
	if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	db, err := sql.Open("sqlite3", "thaichana.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO visits VALUES(?, ?);", chk.ID, chk.PlaceID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "density": "ok" }`))
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}
