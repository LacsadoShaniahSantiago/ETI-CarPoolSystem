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

type Profile struct {
	FirstName      string `json:"First Name"`
	LastName       string `json:"Last Name"`
	Contact        string `json:"Contact"`
	EmailAddr      string `json:"Email Address"`
	AccountCreated string `json:"Account Created"`
	UserType       string `json:"User Type"`
	//Using sql.NullString to allow null values
	LicenseNo  sql.NullString `json:"License No"`
	CarPlateNo sql.NullString `json:"Car Plate No"`
}

var profileList map[string]Profile = map[string]Profile{}

var (
	db          *sql.DB
	err         error
	user        Profile   = Profile{}
	currentTime time.Time = time.Now()
)

func main() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpooling_db")
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/carpool/user/{accountid}", account).Methods("POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/carpool/user", accounts)
	fmt.Println("Listening at port 4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func getAccount() map[string]Profile {

	var users map[string]Profile = map[string]Profile{}

	results, err := db.Query("SELECT * FROM Account")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var u Profile
		var uID string

		err = results.Scan(&uID, &u.FirstName, &u.LastName, &u.Contact, &u.EmailAddr, &u.AccountCreated, &u.UserType, &u.LicenseNo, &u.CarPlateNo)
		if err != nil {
			panic(err.Error())
		}
		users[uID] = u
	}

	return users
}

func accounts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enters Account Function")

	profileList = getAccount()
	userJSON := struct {
		Profiles map[string]Profile `json:"Profiles"`
	}{profileList}

	json.NewEncoder(w).Encode(userJSON)
}

func account(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enters Accounts Function")
	params := mux.Vars(r)

	if r.Method == "POST" {
		fmt.Println("Enters POST Function")
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &user); err == nil {
				if _, ok := userExist(params["accountid"]); !ok {

					fmt.Printf("Inserting Account: %v", user)
					insertAccount(params["accountid"], user)
					w.WriteHeader(http.StatusAccepted)
				} else {

					fmt.Printf("Did not insert account: %v", user)
					fmt.Fprintf(w, "Account ID exists")
					w.WriteHeader(http.StatusConflict)
				}

			} else {
				fmt.Println(err)
			}

		}
	} else if r.Method == "PUT" {
		fmt.Println("Enters PUT Function")
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			if err := json.Unmarshal(body, &user); err == nil {
				if _, ok := userExist(params["accountid"]); ok {
					fmt.Printf("Updating account %v: %v", params["accountid"], user)

					updateAccount(params["accountid"], user)
					w.WriteHeader(http.StatusAccepted)

				} else {
					fmt.Printf("Not updating account: %v", user)
					fmt.Fprintf(w, "Account ID does not exist")
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if user, ok := userExist(params["accountid"]); ok {
		fmt.Println("Enters UserExistCheck Function")

		//Auditable Date
		fmt.Println(user.AccountCreated)
		t, err := time.Parse("2006-01-02 15:04:05", user.AccountCreated)
		if err != nil {
			panic(err)
		}

		t = t.AddDate(1, 0, 0)

		var auditableYear int = t.Year()
		var currentYear int = currentTime.Year()

		fmt.Println("Auditable Year: ", auditableYear)
		fmt.Println("Current Time: ", currentYear)

		var dateDiff = currentYear - auditableYear

		if dateDiff < 0 {
			fmt.Println("Enters Date Diff <= 0 Condition")
			fmt.Fprintln(w, "Account has not passed 1-year data rention.")

		} else if dateDiff >= 1 {
			fmt.Println("Enters Date Diff >= 1 Condition")
			if r.Method == "DELETE" {
				fmt.Println("Enters Delete Function")
				deleteAccount(params["accountid"])
				fmt.Fprintln(w, "Profile "+user.EmailAddr+" deleted")
			} else {
				json.NewEncoder(w).Encode(user)
			}
		}
	} else {
		fmt.Fprintf(w, "Invalid Profile ID")
		w.WriteHeader(http.StatusNotFound)
	}
}

func userExist(id string) (Profile, bool) {
	fmt.Println("Enters UserExist Function")

	var u Profile

	result := db.QueryRow("SELECT * FROM Account WHERE PassengerID=?", id)

	//To nil check result before Scan
	if result != nil {
		err := result.Scan(&id, &u.FirstName, &u.LastName, &u.Contact, &u.EmailAddr, &u.AccountCreated, &u.UserType, &u.LicenseNo, &u.CarPlateNo)
		if err == sql.ErrNoRows {
			return u, false
		}
	}

	return u, true
}

func insertAccount(id string, u Profile) {
	fmt.Println("Enters InsertAccount Function")

	_, err := db.Exec("INSERT INTO Account (PassengerID, FirstName, LastName, Contact, EmailAddr, AccountCreated, UserType, LicenseNo, CarPlateNo) VALUES (?,?,?,?,?,?,?,?,?)", id, u.FirstName, u.LastName, u.Contact, u.EmailAddr, u.AccountCreated, u.UserType, u.LicenseNo.String, u.CarPlateNo.String)
	if err != sql.ErrNoRows {
		panic(err.Error())
	}
}

func updateAccount(id string, u Profile) {
	fmt.Println("Enters UpdateAccount Function")

	_, err := db.Exec("UPDATE Account SET FirstName=?, LastName=?, Contact=?, EmailAddr=?, UserType=?, LicenseNo=?, CarPlateNo=? WHERE PassengerID=?", u.FirstName, u.LastName, u.Contact, u.EmailAddr, u.UserType, u.LicenseNo.String, u.CarPlateNo, id)
	if err != nil {
		panic(err.Error())
	}
}

func deleteAccount(id string) (int64, error) {
	fmt.Println("Enters DeleteAccount Function")

	result, err := db.Exec("DELETE FROM Account WHERE PassengerID=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
