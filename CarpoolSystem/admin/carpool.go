package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Profile struct {
	FirstName      string `json:"First Name"`
	LastName       string `json:"Last Name"`
	Contact        string `json:"Contact"`
	EmailAddr      string `json:"Email Address"`
	AccountCreated string `json:"Account Created"`
	UserType       string `json:"User Type"`

	//NullString struct { String string Valid bool}
	LicenseNo  sql.NullString `json:"License No"`
	CarPlateNo sql.NullString `json:"Car Plate No"`
}

type Profiles struct {
	Profiles map[string]Profile `json:"Profiles"`
}

type Trip struct {
	PassengerID     string `json:"PassengerID"`
	PickUpAddr      string `json:"PickUp Addr"`
	AltPickUpAddr   string `json:"Alt PickUp Addr"`
	StartTravelTime string `json:"Start Travel Time"`
	DestinationAddr string `json:"Destination Addr"`
	PassengerPax    int    `json:"Passenger pax"`
}

type Trips struct {
	Trips map[string]Trip `json:"Trips"`
}

type Enrol struct {
	TripID      string `json:"Trip ID"`
	PassengerID string `json:"Passenger ID"`
	TripStatus  string `json:"Trip Status"`
}

type Enrols struct {
	Enrols map[string]Enrol `json:"Enrols"`
}

var (
	user     Profile
	rider    Trip
	enrolled Enrol

	postUserBody  []byte
	postTripBody  []byte
	postEnrolBody []byte
	currentTime   time.Time = time.Now().Local()
)

func main() {
outer:
	for {
		fmt.Println(strings.Repeat("=", 15))
		fmt.Println("Welcome to Car Pooling Service")
		fmt.Println("1. Create Account")
		fmt.Println("2. Update Account")
		fmt.Println("3. Upgrade Account")
		fmt.Println("4. Delete Account")
		fmt.Println("5. Enrol to Trip")
		fmt.Println("6. View User Enrols")
		fmt.Println("7. View All Trips")
		fmt.Println("8. Create Trip")
		fmt.Println("9. Cancel Trip")
		fmt.Println("")
		fmt.Println("0. Quit")
		fmt.Println("Please press enter to proceed...")
		fmt.Println("")

		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		fmt.Print("Enter an option: ")

		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			createAccount()

		case 2:
			updateAccount()

		case 3:
			upgradeAccount()

		case 4:
			deleteAccount()

		case 5:
			enroltoTrip()
			checkVacancies()

		case 6:
			var uID string = anyUserLogIn()
			viewAllUserEnrols(uID)

		case 7:
			viewAllTrips()

		case 8:
			createTrip()

		case 9:
			cancelTrip()

		case 0:
			break outer

		default:
			fmt.Println("### Enter options available")
		}
	}
}

func profileDisplay(user Profile) {
	if user.UserType == "P" {
		fmt.Println("Passenger Account")
	} else {
		fmt.Println("Car Owner Account")
	}
	fmt.Printf("Full Name\t:  %v %v\n", user.FirstName, user.LastName)
	fmt.Println("Contact\t\t: ", user.Contact)
	fmt.Println("Email Address\t: ", user.EmailAddr)

	//NullTime
	fmt.Println("Account Created\t: ", user.AccountCreated)

	if user.UserType == "D" {
		fmt.Println("-Car Owner Details-")
		fmt.Println("License No	: ", user.LicenseNo.String)
		fmt.Println("Car Plate No	: ", user.CarPlateNo.String)
	}
	fmt.Println()
}

func viewAllAccount() {
	fmt.Println("View All Accounts")
	fmt.Println("")

	client := &http.Client{}

	if getReq, err := http.NewRequest(http.MethodGet, "http://localhost:4000/api/v1/carpool/user", nil); err == nil {
		if res, err := client.Do(getReq); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Profiles
				json.Unmarshal(body, &res)

				for i, u := range res.Profiles {
					fmt.Println("Profile ", i)
					profileDisplay(u)
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func viewPassengerAccount() {
	fmt.Println("View Passenger Accounts")
	fmt.Println("")

	client := &http.Client{}

	if getReq, err := http.NewRequest(http.MethodGet, "http://localhost:4000/api/v1/carpool/user", nil); err == nil {
		if res, err := client.Do(getReq); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Profiles
				json.Unmarshal(body, &res)

				for i, u := range res.Profiles {
					if u.UserType == "P" {
						fmt.Println("Profile ", i)
						profileDisplay(u)
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func viewCarOwnerAccount() {
	fmt.Println("View Car Owner Accounts")
	fmt.Println("")

	client := &http.Client{}

	if getReq, err := http.NewRequest(http.MethodGet, "http://localhost:4000/api/v1/carpool/user", nil); err == nil {
		if res, err := client.Do(getReq); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Profiles
				json.Unmarshal(body, &res)

				for i, u := range res.Profiles {
					if u.UserType == "D" {
						fmt.Println("Profile ", i)
						profileDisplay(u)
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func getAccount(id string) Profile {
	client := &http.Client{}

	if getReq, err := http.NewRequest(http.MethodGet, "http://localhost:4000/api/v1/carpool/user", nil); err == nil {
		if res, err := client.Do(getReq); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Profiles
				json.Unmarshal(body, &res)

				for i, u := range res.Profiles {
					if i == id {
						user = u
						return user
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

	return user
}

func createAccount() {
	var newUser Profile
	var newuserID string

	fmt.Println("Creating A New Account")
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter new account ID\t: ")
	fmt.Scanf("%v", &newuserID)
	reader.ReadString('\n')

	fmt.Print("Enter First name\t: ")
	fmt.Scanf("%v", &(newUser.FirstName))
	reader.ReadString('\n')

	fmt.Print("Enter Last name\t\t: ")
	fmt.Scanf("%v", &(newUser.LastName))
	reader.ReadString('\n')

	fmt.Print("Enter Contact\t\t: ")
	fmt.Scanf("%v", &(newUser.Contact))
	reader.ReadString('\n')

	fmt.Print("Enter Email address\t: ")
	fmt.Scanf("%v", &(newUser.EmailAddr))
	reader.ReadString('\n')

	//Default values
	newUser.AccountCreated = currentTime.Format(time.DateTime)
	newUser.UserType = "P"
	newUser.LicenseNo.Valid = false
	newUser.CarPlateNo.Valid = false

	postUserBody, _ = json.Marshal(newUser)
	resBody := bytes.NewBuffer(postUserBody)

	client := &http.Client{}
	if postReq, err := http.NewRequest(http.MethodPost, "http://localhost:4000/api/v1/carpool/user/"+newuserID, resBody); err == nil {
		if res, err := client.Do(postReq); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Profile ", newuserID, " created successfully")
				profileDisplay(newUser)

			} else if res.StatusCode == 409 {
				fmt.Println("Error - Profile ", newuserID, "exists")
			}
		} else {
			fmt.Println("Profile ", newuserID, " created successfully")
			profileDisplay(newUser)
			//fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func updateAccount() {
	var updateAccount Profile
	var updateID string

	viewAllAccount()
	fmt.Println(strings.Repeat("-", 15))
	fmt.Println("Updating Account")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter account ID to update\t: ")
	fmt.Scanf("%v", &updateID)
	reader.ReadString('\n')

	updateAccount = getAccount(updateID)

	fmt.Printf("Enter updated First name\t: ")
	fmt.Scanf("%v", &(updateAccount.FirstName))
	reader.ReadString('\n')

	fmt.Printf("Enter updated Last name\t\t: ")
	fmt.Scanf("%v", &(updateAccount.LastName))
	reader.ReadString('\n')

	fmt.Printf("Enter updated Contact\t\t: ")
	fmt.Scanf("%v", &(updateAccount.Contact))
	reader.ReadString('\n')

	fmt.Printf("Enter updated Email Address\t: ")
	fmt.Scanf("%v", &(updateAccount.EmailAddr))
	reader.ReadString('\n')

	if updateAccount.UserType == "D" {
		fmt.Printf("Enter updated License No\t: ")
		fmt.Scanf("%v", &(updateAccount.LicenseNo.String))
		reader.ReadString('\n')

		fmt.Printf("Enter updated License No\t: ")
		fmt.Scanf("%v", &(updateAccount.CarPlateNo.String))
		reader.ReadString('\n')
	}

	fmt.Println("")

	postUserBody, _ := json.Marshal(updateAccount)

	client := &http.Client{}
	if putReq, err := http.NewRequest(http.MethodPut, "http://localhost:4000/api/v1/carpool/user/"+updateID, bytes.NewBuffer(postUserBody)); err == nil {
		var updatedAccount Profile
		if res, err := client.Do(putReq); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Account ", updateID, " updated successfully")

				updatedAccount = getAccount(updateID)
				profileDisplay(updatedAccount)

			} else if res.StatusCode == 409 {
				fmt.Println("Account has not passed 1-year data rention.")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func upgradeAccount() {
	var upgradeID string
	var upgradingAccount Profile

	viewPassengerAccount()
	fmt.Println(strings.Repeat("-", 15))
	fmt.Println("Upgrading Account")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter account ID to upgrade\t: ")
	fmt.Scanf("%v", &upgradeID)
	reader.ReadString('\n')

	upgradingAccount = getAccount(upgradeID)

	fmt.Print("Enter License No\t\t: ")
	fmt.Scanf("%v", &(upgradingAccount.LicenseNo.String))
	reader.ReadString('\n')

	fmt.Print("Enter Car Plate No\t\t: ")
	fmt.Scanf("%v", &(upgradingAccount.CarPlateNo.String))
	reader.ReadString('\n')

	//Update UserType P to D
	upgradingAccount.UserType = "D"

	//Set NullString Valid = True
	upgradingAccount.LicenseNo.Valid = true
	upgradingAccount.CarPlateNo.Valid = true

	postUserBody, _ := json.Marshal(upgradingAccount)

	client := &http.Client{}
	if putReq, err := http.NewRequest(http.MethodPut, "http://localhost:4000/api/v1/carpool/user/"+upgradeID, bytes.NewBuffer(postUserBody)); err == nil {
		var upgradedAccount Profile
		if res, err := client.Do(putReq); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("")
				fmt.Println("Profile ", upgradeID, " upgraded successfully")

				upgradedAccount = getAccount(upgradeID)
				profileDisplay(upgradedAccount)

			} else if res.StatusCode == 404 {
				fmt.Println("Error - Profile", upgradeID, " does not exist")
			}
		} else {
			//fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func deleteAccount() {
	var deleteID string

	viewAllAccount()

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter account ID to be deleted: ")
	fmt.Scanf("%v", &deleteID)
	reader.ReadString('\n')

	//Evaluate Date Difference
	user = getAccount(deleteID)
	t, err := time.Parse("2006-01-02 15:04:05", user.AccountCreated)
	if err != nil {
		panic(err)
	}
	t = t.AddDate(1, 0, 0)
	var auditableYear int = t.Year()
	var currentYear int = currentTime.Year()
	var datediff = currentYear - auditableYear

	client := &http.Client{}
	if delReq, err := http.NewRequest(http.MethodDelete, "http://localhost:4000/api/v1/carpool/user/"+deleteID, nil); err == nil {
		if res, err := client.Do(delReq); err == nil {
			if res.StatusCode == 202 || datediff >= 1 {
				fmt.Println("Profile ", deleteID, " deleted successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Profile ", deleteID, " does not exist")
			} else if res.StatusCode == 406 || datediff < 0 {
				fmt.Println("Account has not passed 1-year data rention.")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func tripDisplay(t Trip) {
	fmt.Println("To ", t.DestinationAddr)
	fmt.Println(strings.Repeat("- ", 10))
	fmt.Println("Pick Up Address\t\t\t: ", t.PickUpAddr)
	fmt.Println("Alternate Pick Up Address\t: ", t.AltPickUpAddr)
	fmt.Println("Schedule Start Time\t\t: ", t.StartTravelTime)
	fmt.Println("Passenger Pax: \t\t\t: ", t.PassengerPax)
	fmt.Println("")
}

func viewAllTrips() {
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/carpool/trip", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Trips
				json.Unmarshal(body, &res)

				fmt.Println("View All Trips")
				fmt.Println("")
				for i, t := range res.Trips {
					fmt.Println("Trip ", i)
					tripDisplay(t)
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func viewUserPublishedTrips(uID string) {
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/carpool/trip", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Trips
				json.Unmarshal(body, &res)

				fmt.Println("View All User Created Trips")
				fmt.Println("")
				for i, t := range res.Trips {
					if t.PassengerID == uID {
						fmt.Println("Trip ", i)
						tripDisplay(t)
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func userLogIn() string {
	var logID string

	fmt.Print("")
	viewCarOwnerAccount()
	fmt.Print("")
	//Account to Log In
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter car owner account id to log in: ")
	fmt.Scanf("%v", &logID)
	reader.ReadString('\n')

	//Retrieve Log Account
	user = getAccount(logID)
	fmt.Println(strings.Repeat("-", 10))
	fmt.Println("Logged In Account")
	profileDisplay(user)
	return logID
}

func anyUserLogIn() string {
	var logID string

	fmt.Print("")
	viewAllAccount()
	fmt.Print("")
	//Account to Log In
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Print("Enter account id to log in: ")
	fmt.Scanf("%v", &logID)
	reader.ReadString('\n')

	//Retrieve Log Account
	user = getAccount(logID)
	fmt.Println(strings.Repeat("-", 10))
	fmt.Println("Logged In Account")
	profileDisplay(user)
	return logID
}

func getTrip(tID string) Trip {
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/carpool/trip", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Trips
				json.Unmarshal(body, &res)
				for i, t := range res.Trips {
					if i == tID {
						rider = t
						return rider
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

	return rider
}

func createTrip() {
	var logID string
	var newTrip Trip
	var newtID string
	var travelTime string

	fmt.Println("Creating Trip")
	fmt.Println("")
	//Identify user logged in
	logID = userLogIn()

	fmt.Println("")
	fmt.Println("Create Trip")
	//fmt.Println("Please press enter to proceed...")

	//Create Trip
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new Trip ID\t\t: ")
	fmt.Scanf("%v", &newtID)
	reader.ReadString('\n')

	fmt.Print("Enter Pickup Address\t\t: ")
	fmt.Scanf("%v", &(newTrip.PickUpAddr))
	reader.ReadString('\n')

	fmt.Print("Enter Alternate Pickup Address\t: ")
	fmt.Scanf("%v", &(newTrip.AltPickUpAddr))
	reader.ReadString('\n')

	fmt.Print("Enter start time (hh:mm:ss)\t: ")
	fmt.Scanf("%v", &travelTime)
	reader.ReadString('\n')

	//Format today's time with time start
	var dateNow = currentTime.Format("2006-01-02")
	var createSchedule = dateNow + " " + travelTime
	newTrip.StartTravelTime = createSchedule
	fmt.Println("Scheduled Time: ", newTrip.StartTravelTime)

	fmt.Print("Enter Desination Address\t: ")
	fmt.Scanf("%v", &(newTrip.DestinationAddr))
	reader.ReadString('\n')

	fmt.Print("Enter Passenger Pax\t\t: ")
	fmt.Scanf("%d", &(newTrip.PassengerPax))
	reader.ReadString('\n')

	//Assign Car Owner Profile ID
	newTrip.PassengerID = logID

	postTripBody, _ := json.Marshal(newTrip)
	resBody := bytes.NewBuffer(postTripBody)

	client := &http.Client{}
	if postReq, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/trip/"+newtID, resBody); err == nil {
		if res, err := client.Do(postReq); err == nil {
			fmt.Println("")
			if res.StatusCode == 202 || res.StatusCode == 200 {
				fmt.Println("Trip ", newtID, " created successfully")
				tripDisplay(newTrip)

			} else if res.StatusCode == 409 {
				fmt.Println("Error - Trip ", newtID, " exists")
			}
		} else {
			fmt.Println("Trip ", newtID, " created successfully")
			tripDisplay(newTrip)
			//fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func cancelTrip() {
	var uID string
	var deleteTID string

	fmt.Print("")
	fmt.Println("Cancelling Trip")
	fmt.Print("")

	//Identify Car Owner Account
	uID = userLogIn()

	viewUserPublishedTrips(uID)

	fmt.Println("")
	fmt.Println("Cancel Trip")
	//Cancel Trip
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Trip ID to Cancel: ")
	fmt.Scanf("%v", &deleteTID)
	reader.ReadString('\n')

	//Evaluate Minute Difference
	rider = getTrip(deleteTID)
	startTime, err := time.Parse("2006-01-02 15:04:05", rider.StartTravelTime)
	if err != nil {
		panic(err)
	}

	timeNow := currentTime

	diff := timeNow.Sub(startTime)
	diffInMinutes := int(diff.Minutes())

	client := &http.Client{}
	if delReq, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/carpool/trip/"+deleteTID, nil); err == nil {
		if res, err := client.Do(delReq); err == nil {
			if res.StatusCode == 200 || diffInMinutes < -30 {
				fmt.Println("Trip ", deleteTID, " deleted successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Trip ", deleteTID, " does not exist")
			} else if res.StatusCode == 406 || diffInMinutes >= -30 {
				fmt.Println("Trip is within 30 minutes before Start Time. No cancellation is allowed")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func enrolDisplay(e Enrol) {
	user = getAccount(e.PassengerID)
	rider = getTrip(e.TripID)

	fmt.Printf("Passenger name\t: %v %v\tAccountID: %v\n", user.FirstName, user.LastName, e.PassengerID)
	fmt.Println(strings.Repeat("- ", 15))
	fmt.Println("Trip From\t: ", rider.PickUpAddr)
	fmt.Println("\tAlt From: ", rider.AltPickUpAddr)
	fmt.Println("Trip To\t\t: ", rider.DestinationAddr)
	fmt.Println("Scehdule Time\t: ", rider.StartTravelTime)
	if e.TripStatus == "F" {
		fmt.Println("Trip Status: Full")
	} else {
		fmt.Println("Trip Status: Vacant")
	}
}

func viewAllUserEnrols(uID string) {
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:6000/api/v1/carpool/enrol", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Enrols
				json.Unmarshal(body, &res)

				fmt.Println("View Passenger's Enrolled Trips")
				fmt.Println("")
				for i, e := range res.Enrols {
					if e.PassengerID == uID {
						fmt.Println("Enrolled Trip ", i)
						enrolDisplay(e)
						fmt.Print("")
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func viewVacantEnrol() {
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:6000/api/v1/carpool/enrol", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Enrols
				json.Unmarshal(body, &res)

				fmt.Println("View Vacant Trips")
				fmt.Println("")
				for _, e := range res.Enrols {
					if e.TripStatus == "V" {
						rider = getTrip(e.TripID)
						tripDisplay(rider)
						fmt.Print("")
					}
				}
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func enrolTrip() string {
	var rideID string

	fmt.Print("")
	viewAllTrips()
	fmt.Print("")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select Trip to Enrol: ")
	fmt.Scanf("%v", &rideID)
	reader.ReadString('\n')

	//Get Trip Details
	rider = getTrip(rideID)
	tripDisplay(rider)

	return rideID
}

func enroltoTrip() {
	var logID string
	var rideID string
	var newEnrol Enrol
	var newEID string

	fmt.Println("Enrolling to Trip")
	fmt.Println("")

	//Identify user logged in
	logID = anyUserLogIn()

	//Identify Trip to enrol
	rideID = enrolTrip()

	//Create Enrol
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new Enrol ID\t: ")
	fmt.Scanf("%v", &newEID)
	reader.ReadString('\n')

	newEnrol.PassengerID = logID
	newEnrol.TripID = rideID
	newEnrol.TripStatus = "V"

	//Display Information of Profile and Ride enrolling
	user = getAccount(logID)
	rider = getTrip(rideID)

	fmt.Print("")
	fmt.Printf("Passenger Enrolling\t: %v %v\n\n", user.FirstName, user.LastName)
	fmt.Println("--Enrolling To Trip--")
	tripDisplay(rider)

	postEnrolBody, _ := json.Marshal(newEnrol)
	resBody := bytes.NewBuffer(postEnrolBody)

	client := &http.Client{}
	if postReq, err := http.NewRequest(http.MethodPut, "http://localhost:6000/api/v1/carpool/enrol/"+newEID, resBody); err == nil {
		if res, err := client.Do(postReq); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Enrol ", newEID, " created successfully")
				enrolDisplay(newEnrol)
			} else if res.StatusCode == 406 {
				fmt.Print("Trip is not vacant, enrol to another trip.")
				viewVacantEnrol()
			} else if res.StatusCode == 409 {
				fmt.Println("Error - Enrol ", newEID, " exists")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

func getNumberOfEnrols(tID string) int {

	var enrolsList map[string]Enrol = map[string]Enrol{}

	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:6000/api/v1/carpool/enrol", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Enrols

				json.Unmarshal(body, &res)

				for i, e := range res.Enrols {
					if e.TripID == tID {
						enrolsList[i] = e
					}
				}

				return len(enrolsList)
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

	return len(enrolsList)
}

func getAllTrips() map[string]Trip {

	var tripsList map[string]Trip = map[string]Trip{}
	client := &http.Client{}

	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5000/api/v1/carpool/trip", nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				var res Trips
				json.Unmarshal(body, &res)

				for i, t := range res.Trips {
					tripsList[i] = t
				}

				return tripsList
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}

	return tripsList
}
func checkVacancies() {

	var trips map[string]Trip = getAllTrips()

	for i, t := range trips {

		var totalEnrols int = getNumberOfEnrols(i)
		if t.PassengerPax == totalEnrols {
			updateEnrol(i)
		}
	}

}

func updateEnrol(eID string) {

	var updateEnrol Enrol
	updateEnrol.TripStatus = "F"

	postEnrolBody, _ := json.Marshal(updateEnrol)

	client := &http.Client{}
	if putReq, err := http.NewRequest(http.MethodPut, "http://localhost:6000/api/v1/enrol/"+eID, bytes.NewBuffer(postEnrolBody)); err == nil {
		if res, err := client.Do(putReq); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Enrol ", eID, " updated successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - enrol ", eID, " does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}
