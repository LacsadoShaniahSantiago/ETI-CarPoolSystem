package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Enrol struct {
	TripID      string `json:"Trip ID"`
	PassengerID string `json:"Passenger ID"`
	TripStatus  string `json:"Trip Status"`
}

var enrolList map[string]Enrol = map[string]Enrol{}

var (
	db        *sql.DB
	err       error
	enrolTrip Enrol = Enrol{}
)

func main() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpooling_db")
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/carpool/enrol/{enrolid}", enrol).Methods("POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/carpool/enrol", enrols)
	fmt.Println("Listening at port 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}

func getEnrol() map[string]Enrol {
	var enrols map[string]Enrol = map[string]Enrol{}

	results, err := db.Query("SELECT * FROM Enrol")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var e Enrol
		var eID string

		err = results.Scan(&eID, &e.TripID, &e.PassengerID, &e.TripStatus)
		if err != nil {
			panic(err.Error())
		}

		enrols[eID] = e
	}

	return enrols
}

func enrols(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enters Enrols Function")

	enrolList = getEnrol()
	enrolJSON := struct {
		Enrols map[string]Enrol `json:"Enrols"`
	}{enrolList}

	json.NewEncoder(w).Encode(enrolJSON)
}

func enrol(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" {
		fmt.Println("Enters POST Function")
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &enrolTrip); err == nil {
				if _, ok := enrolExist(params["enrolid"]); !ok {

					if conflict := checkConflict(enrolTrip); !conflict {
						fmt.Println("Inserting Enrol: ", enrolTrip)
						insertEnrol(params["enrolid"], enrolTrip)
						w.WriteHeader(http.StatusAccepted)
					} else {
						fmt.Println("Enrolment is not vacant")
						w.WriteHeader(http.StatusNotAcceptable)
					}
				} else {
					fmt.Printf("Did not insert enrol: %v", enrolTrip)
					fmt.Fprintf(w, "Enrol ID exists")
					w.WriteHeader(http.StatusConflict)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		fmt.Println("Enters PUT Function")

		if body, err := ioutil.ReadAll(r.Body); err == nil {
			if err := json.Unmarshal(body, &enrolTrip); err == nil {
				if _, ok := enrolExist(params["enrolid"]); ok {
					fmt.Printf("Updating enrol: ", enrolTrip)

					updateStatus(params["enrolid"], enrolTrip)
					w.WriteHeader(http.StatusAccepted)

				} else {
					fmt.Print("Not updating enrolment: ", enrolTrip)
					fmt.Fprintf(w, "Enrolment ID does not exist")
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}

func enrolExist(id string) (Enrol, bool) {
	fmt.Println("Enters EnrolExist Function")

	var e Enrol

	result := db.QueryRow("SELECT * FROM Enrol WHERE EnrolID=?", id)
	err := result.Scan(&id, &e.TripID, &e.PassengerID, &e.TripStatus)
	if err == sql.ErrNoRows {
		return e, false
	}

	return e, true
}

func checkConflict(e Enrol) bool {
	var Enrols map[string]Enrol = getEnrol()

	for _, enrol := range Enrols {

		//Check Users Enrol No Conflict
		if enrol.PassengerID == e.PassengerID {
			if enrol.TripID == e.TripID || e.TripStatus == "F" {
				return true
			}
		}
	}

	return false
}

func insertEnrol(id string, e Enrol) {
	fmt.Println("Enters InsertEnrol Function")

	_, err := db.Exec("INSERT INTO Enrol (EnrolID, TripID, PassengerID, TripStatus) VALUES (?,?,?,?)", id, e.TripID, e.PassengerID, e.TripStatus)
	if err != sql.ErrNoRows {
		panic(err.Error())
	}
}

func updateStatus(id string, e Enrol) {
	fmt.Println("Enters UpdateStatus Function")

	_, err := db.Exec("UPDATE Enrol SET TripStatus=? WHERE TripID=? ", e.TripStatus, id)
	if err != nil {
		panic(err.Error())
	}
}
