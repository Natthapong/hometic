package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pair struct {
	DeviceID int64
	UserID   int64
}

type PairDeviceHandler struct {
	createPairDevice CreatePairDevice
}

type CreatePairDevice func(p Pair) error

var createPairDevice = func(p Pair) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error: ", err)
	}
	_, err = db.Exec("INSERT INTO pairs VALUES ($1,$2);", p.DeviceID, p.UserID)
	return err
}

func (ph *PairDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var p Pair
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()
	fmt.Printf("pair: %#v\n", p)

	err = ph.createPairDevice(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.Write([]byte(`{"status":"active"}`))
}

func setupRouter() {
	r := mux.NewRouter()
	r.Handle("/pair-device", &PairDeviceHandler{createPairDevice}).Methods(http.MethodPost)

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("...Starting...")
	log.Fatal(server.ListenAndServe()) // listen and serve on 127.0.0.0:8080
}

func main() {
	setupRouter()
}
