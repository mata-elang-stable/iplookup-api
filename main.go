package main

import (
	"encoding/json"
	models2 "github.com/fadhilyori/iplookup-go/internal/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Env for holding connection
type Env struct {
	reader *models2.IPLookup
}

func main() {
	//mmdbPath := os.Args[1]
	//
	//if mmdbPath == "" {
	//	log.Panic("Usage: ./main /path/to/mmdb_file <port>\nport default: 80")
	//}

	log.Print("Starting up service ...")
	//db, err := models.NewReaderHandler(mmdbPath)
	//if err != nil {
	//	log.Panic(err)
	//}
	//

	db := models2.NewIPLookup(
		"./assets/GeoLite2-City.mmdb",
		"./assets/GeoLite2-ASN.mmdb",
	)

	env := &Env{db}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{ipaddress}", env.lookUp).Methods("GET")
	log.Print("Service running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func (env *Env) lookUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	vars := mux.Vars(r)
	ipaddress := vars["ipaddress"]

	res := env.reader.Lookup(ipaddress)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("Error encoding json: %v\n", err)
		return
	}
}
