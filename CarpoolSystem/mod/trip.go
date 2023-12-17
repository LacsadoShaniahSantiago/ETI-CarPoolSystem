package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Trip struct {
	PassengerID     string `json:"PassengerID"`
	PickUpAddr      string `json:"PickUp Addr"`
	AltPickUpAddr   string `json:"Alt PickUp Addr"`
	StartTravelTime string `json:"Start Travel Time"`
	DestinationAddr string `json:"Destination Addr"`
	PassengerPax    int    `json:"Passenger pax"`
}

var tripList map[string]Trip = map[string]Trip{}

var (
	db          *sql.DB
	err         error
	ride        Trip      = Trip{}
	currentTime time.Time = time.Now()
)

func main() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpooling_db")
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/carpool/trip/{tripid}", trip).Methods("POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/carpool/trip", trips)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}

func getTrip() map[string]Trip {

	var trips map[string]Trip = map[string]Trip{}

	results, err := db.Query("SELECT * FROM Trip")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var t Trip
		var tID string

		err = results.Scan(&tID, &t.PassengerID, &t.PickUpAddr, &t.AltPickUpAddr, &t.StartTravelTime, &t.DestinationAddr, &t.PassengerPax)
		if err != nil {
			panic(err.Error())
		}

		trips[tID] = t
	}

	return trips
}

func trips(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enters Trips Function")

	tripList = getTrip()
	tripJSON := struct {
		Trips map[string]Trip `json:"Trips"`
	}{tripList}

	json.NewEncoder(w).Encode(tripJSON)
}

func trip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &ride); err == nil {
				if _, ok := tripExist(params["tripid"]); !ok {

					fmt.Println("Inserting Trip: ", ride)
					insertTrip(params["tripid"], ride)
					w.WriteHeader(http.StatusAccepted)

				} else {
					fmt.Printf("Did not insert trip: %v", ride)
					fmt.Fprintf(w, "Trip ID exists")
					w.WriteHeader(http.StatusConflict)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if ride, ok := tripExist(params["tripid"]); ok {
		fmt.Println("Enter TripExistCheck Function")

		//Delete Timing
		startTime, err := time.Parse("2006-01-02 15:04:05", ride.StartTravelTime)
		if err != nil {
			panic(err)
		}

		timeNow := currentTime

		diff := timeNow.Sub(startTime)
		diffInMinutes := int(diff.Minutes())

		fmt.Println("Difference in Minutes: ", diffInMinutes)

		//When DiffInMinutes less than -30, Means 30 Mins More Before Starts Time
		if diffInMinutes > -30 {
			fmt.Println("Enters Minute Diff > -30 Condition")
			if r.Method == "DELETE" {
				fmt.Println("Enters DELETE Function")
				deleteTrip(params["tripid"])
				fmt.Fprintln(w, "Trip "+params["tripid"], " deleted")
			} else {
				json.NewEncoder(w).Encode(ride)
			}

		} else if diffInMinutes <= -30 {
			fmt.Println("Enters Minute Diff <= -30 Condition")
			fmt.Fprintln(w, "Trip within 30 minutes before Start Time. No cancellation is allowed")
		}
	} else {
		fmt.Fprintf(w, "Invalid Trip ID")
		w.WriteHeader(http.StatusNotFound)
	}
}

func tripExist(id string) (Trip, bool) {
	fmt.Println("Enters TripExist Function")

	var t Trip

	result := db.QueryRow("SELECT * FROM Trip WHERE TripID=?", id)

	//To nil check result before Scan
	if result != nil {
		err := result.Scan(&id, &t.PassengerID, &t.PickUpAddr, &t.AltPickUpAddr, &t.StartTravelTime, &t.DestinationAddr, &t.PassengerPax)
		if err == sql.ErrNoRows {
			return t, false
		}
	}

	return t, true
}

func insertTrip(id string, t Trip) {
	fmt.Println("Enters InsertTrip Function")

	_, err := db.Exec("INSERT INTO Trip (TripID, PassengerID, PickUpAddr, AltpickUpAddr, StartTravelTime, DestinationAddr, PassengerPax) VALUES (?,?,?,?,?,?,?)", id, t.PassengerID, t.PickUpAddr, t.AltPickUpAddr, t.StartTravelTime, t.DestinationAddr, t.PassengerPax)
	if err != sql.ErrNoRows {
		panic(err.Error())
	}
}

func deleteTrip(id string) (int64, error) {
	fmt.Println("Enters DeleteTrip Function")

	result, err := db.Exec("DELETE FROM Trip WHERE TripID=?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func deleteUserTrips(uID string) (int64, error) {
	fmt.Println("Enters Delete All User's Trip Function")

	result, err := db.Exec("DELETE FROM Trip WHERE PassengerID=?", uID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
